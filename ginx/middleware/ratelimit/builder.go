package ratelimit

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"pkg404/limiter"

	"github.com/gin-gonic/gin"
)

type Builder struct {
	prefix  string
	limiter limiter.Limiter
}

func NewBuilder(l limiter.Limiter) *Builder {
	return &Builder{
		prefix:  "ip-limiter",
		limiter: l,
	}
}

func (b *Builder) Prefix(prefix string) *Builder {
	b.prefix = prefix
	return b
}

func (b *Builder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.GetHeader("x-stress") == "true" {
			newCtx := context.WithValue(ctx, "x-stress", true)
			ctx.Request = ctx.Request.Clone(newCtx)
			ctx.Next()
			return
		}
		limited, err := b.limiter.Limit(ctx, fmt.Sprintf("%s:%s", b.prefix, ctx.ClientIP()))
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if limited {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		ctx.Next()
	}
}
