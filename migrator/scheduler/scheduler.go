package scheduler

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/to404hanga/pkg404/ginx"
	"github.com/to404hanga/pkg404/gormx/connpool"
	"github.com/to404hanga/pkg404/logger"
	"github.com/to404hanga/pkg404/migrator"
	"github.com/to404hanga/pkg404/migrator/events"
	"github.com/to404hanga/pkg404/migrator/validator"
	"gorm.io/gorm"
)

// Scheduler 统一管理整个迁移过程
type Scheduler[T migrator.Entity] struct {
	lock       sync.Mutex
	src        *gorm.DB
	dst        *gorm.DB
	pool       *connpool.DoubleWritePool
	l          logger.Logger
	pattern    string
	cancelFull func()
	cancelIncr func()
	producer   events.Producer
	fulls      map[string]func() // 如果允许多个全量校验同时运行
}

func NewScheduler[T migrator.Entity](l logger.Logger, src, dst *gorm.DB, pool *connpool.DoubleWritePool, producer events.Producer) *Scheduler[T] {
	return &Scheduler[T]{
		lock:       sync.Mutex{},
		src:        src,
		dst:        dst,
		pool:       pool,
		l:          l,
		pattern:    "*",
		cancelFull: func() {},
		cancelIncr: func() {},
		producer:   producer,
		fulls:      make(map[string]func()),
	}
}

func (s *Scheduler[T]) RegisterRoutes(server *gin.RouterGroup) {
	server.POST("/src_only", ginx.Wrap(s.SrcOnly))
	server.POST("/src_first", ginx.Wrap(s.SrcFirst))
	server.POST("/dst_only", ginx.Wrap(s.DstOnly))
	server.POST("/dst_first", ginx.Wrap(s.DstFirst))
	server.POST("/full/start", ginx.Wrap(s.StartFullValidation))
	server.POST("/full/stop", ginx.Wrap(s.StopFullValidation))
	server.POST("/incr/start", ginx.WrapBody[StartIncrRequest](s.StartIncrValidation))
	server.POST("/incr/stop", ginx.Wrap(s.StopIncrValidation))
}

func (s *Scheduler[T]) SrcOnly(ctx *gin.Context) (ginx.Result, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.pattern = connpool.PatternSrcOnly
	s.pool.UpdatePattern(s.pattern)
	return ginx.Result{
		Code: http.StatusOK,
		Msg:  "OK",
	}, nil
}

func (s *Scheduler[T]) SrcFirst(ctx *gin.Context) (ginx.Result, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.pattern = connpool.PatternSrcFirst
	s.pool.UpdatePattern(s.pattern)
	return ginx.Result{
		Code: http.StatusOK,
		Msg:  "OK",
	}, nil
}

func (s *Scheduler[T]) DstOnly(ctx *gin.Context) (ginx.Result, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.pattern = connpool.PatternDstOnly
	s.pool.UpdatePattern(s.pattern)
	return ginx.Result{
		Code: http.StatusOK,
		Msg:  "OK",
	}, nil
}

func (s *Scheduler[T]) DstFirst(ctx *gin.Context) (ginx.Result, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.pattern = connpool.PatternDstFirst
	s.pool.UpdatePattern(s.pattern)
	return ginx.Result{
		Code: http.StatusOK,
		Msg:  "OK",
	}, nil
}

func (s *Scheduler[T]) StopIncrValidation(ctx *gin.Context) (ginx.Result, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.cancelIncr()
	return ginx.Result{
		Code: http.StatusOK,
		Msg:  "OK",
	}, nil
}

func (s *Scheduler[T]) StartIncrValidation(ctx *gin.Context, req StartIncrRequest) (ginx.Result, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	cancel := s.cancelIncr
	v, err := s.newValidator()
	if err != nil {
		s.l.Error("创建增量校验器失败", logger.Error(err))
		return ginx.Result{
			Code: http.StatusInternalServerError,
			Msg:  "系统异常",
		}, err
	}
	v.UpdateTime(req.UpdateTime).SleepInterval(time.Duration(req.Interval) * time.Millisecond)
	var goCtx context.Context
	goCtx, s.cancelIncr = context.WithCancel(context.Background())

	go func() {
		cancel()
		err := v.Validate(goCtx)
		s.l.Warn("退出增量校验", logger.Error(err))
	}()

	return ginx.Result{
		Code: http.StatusOK,
		Msg:  "启动增量校验成功",
	}, nil
}

func (s *Scheduler[T]) StopFullValidation(ctx *gin.Context) (ginx.Result, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.cancelFull()
	return ginx.Result{
		Code: http.StatusOK,
		Msg:  "OK",
	}, nil
}

func (s *Scheduler[T]) StartFullValidation(ctx *gin.Context) (ginx.Result, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	cancel := s.cancelFull
	v, err := s.newValidator()
	if err != nil {
		s.l.Error("创建全量校验器失败", logger.Error(err))
		return ginx.Result{
			Code: http.StatusInternalServerError,
			Msg:  "系统异常",
		}, err
	}
	var goCtx context.Context
	goCtx, s.cancelFull = context.WithCancel(context.Background())

	go func() {
		cancel()
		err := v.Validate(goCtx)
		s.l.Warn("退出全量校验", logger.Error(err))
	}()
	return ginx.Result{
		Code: http.StatusOK,
		Msg:  "启动全量校验成功",
	}, nil
}

func (s *Scheduler[T]) newValidator() (*validator.Validator[T], error) {
	switch s.pattern {
	case connpool.PatternSrcOnly, connpool.PatternSrcFirst:
		return validator.NewValidator[T](s.src, s.dst, "SRC", s.l, s.producer), nil
	case connpool.PatternDstOnly, connpool.PatternDstFirst:
		return validator.NewValidator[T](s.dst, s.src, "DST", s.l, s.producer), nil
	default:
		return nil, fmt.Errorf("未知的 pattern %s", s.pattern)
	}
}

type StartIncrRequest struct {
	UpdateTime int64 `json:"update_time"`
	Interval   int64 `json:"interval"`
}
