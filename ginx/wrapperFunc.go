package ginx

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/to404hanga/pkg404/logger"
)

var (
	L      logger.Logger = logger.NewNopLogger()
	vector *prometheus.CounterVec
)

func InitCounter(opt prometheus.CounterOpts) {
	vector = prometheus.NewCounterVec(opt, []string{"code"})
	prometheus.MustRegister(vector)
}

func WrapBodyAndClaims[Req interface{}, Claims interface{}](bizFunc func(ctx *gin.Context, req Req, claims Claims) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.Bind(&req); err != nil {
			L.Error("输入错误", logger.Error(err))
			return
		}
		L.Debug("输入参数", logger.Any("req", req))
		val, ok := ctx.Get("user")
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		claims, ok := val.(Claims)
		if !ok {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		res, err := bizFunc(ctx, req, claims)
		vector.WithLabelValues(strconv.Itoa(res.Code)).Inc()
		if err != nil {
			L.Error("执行业务逻辑失败", logger.Error(err))
		}
		ctx.JSON(http.StatusOK, res)
	}
}

func WrapBody[Req interface{}](bizFunc func(ctx *gin.Context, req Req) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.Bind(&req); err != nil {
			L.Error("输入错误", logger.Error(err))
			return
		}
		L.Debug("输入参数", logger.Any("req", req))
		res, err := bizFunc(ctx, req)
		vector.WithLabelValues(strconv.Itoa(res.Code)).Inc()
		if err != nil {
			L.Error("执行业务逻辑失败", logger.Error(err))
		}
		ctx.JSON(http.StatusOK, res)
	}
}

func WrapClaims[Claims interface{}](bizFunc func(ctx *gin.Context, claims Claims) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		val, ok := ctx.Get("user")
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		claims, ok := val.(Claims)
		if !ok {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		res, err := bizFunc(ctx, claims)
		vector.WithLabelValues(strconv.Itoa(res.Code)).Inc()
		if err != nil {
			L.Error("执行业务逻辑失败", logger.Error(err))
		}
		ctx.JSON(http.StatusOK, res)
	}
}

func Wrap(bizFunc func(ctx *gin.Context) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, err := bizFunc(ctx)
		vector.WithLabelValues(strconv.Itoa(res.Code)).Inc()
		if err != nil {
			L.Error("执行业务逻辑失败", logger.String("path", ctx.Request.URL.Path), logger.String("route", ctx.FullPath()), logger.Error(err))
		}
		ctx.JSON(http.StatusOK, res)
	}
}
