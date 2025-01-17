package grpcx

import (
	"context"
	"net"
	"strconv"
	"time"

	"github.com/LinSrHang/pkg404/logger"
	"github.com/LinSrHang/pkg404/netx"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"google.golang.org/grpc"
)

type Server struct {
	*grpc.Server
	Port        int
	EtcdTTL     int64
	EtcdClient  *clientv3.Client
	etcdManager endpoints.Manager
	etcdKey     string
	cancel      func()
	Name        string
	L           logger.Logger
}

// Serve 启动服务器并阻塞
func (s *Server) Serve() error {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	port := strconv.Itoa(s.Port)
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	err = s.register(ctx, port)
	if err != nil {
		return err
	}
	return s.Server.Serve(l)
}

func (s *Server) register(ctx context.Context, port string) error {
	cli := s.EtcdClient
	serviceName := "service/" + s.Name
	em, err := endpoints.NewManager(cli, serviceName)
	if err != nil {
		return err
	}
	s.etcdManager = em
	ip := netx.GetOutboundIP()
	s.etcdKey = serviceName + "/" + ip
	addr := ip + ":" + port
	leaseResp, err := cli.Grant(ctx, s.EtcdTTL)
	if err != nil {
		return err
	}
	// 开启续约
	ch, err := cli.KeepAlive(ctx, leaseResp.ID)
	if err != nil {
		return err
	}
	go func() {
		for chResp := range ch {
			s.L.Debug("续约: ", logger.String("resp", chResp.String()))
		}
	}()
	return em.AddEndpoint(ctx, s.etcdKey, endpoints.Endpoint{Addr: addr}, clientv3.WithLease(leaseResp.ID))
}

func (s *Server) Close() error {
	s.cancel()
	if s.etcdManager != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := s.etcdManager.DeleteEndpoint(ctx, s.etcdKey)
		if err != nil {
			return err
		}
	}
	err := s.EtcdClient.Close()
	if err != nil {
		return err
	}
	s.Server.GracefulStop()
	return nil
}
