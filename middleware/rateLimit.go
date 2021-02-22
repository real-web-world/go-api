package middleware

import (
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"golang.org/x/time/rate"

	"github.com/real-web-world/go-api/pkg/gin"
)

var (
	limiterSet     = cache.New(5*time.Minute, 10*time.Minute)
	getReqClientIP = func(c *gin.Context) string {
		return c.ClientIP()
	}
	getLimitFn = func(reqDelay time.Duration, limit int) func(*gin.Context) (*rate.Limiter,
		time.Duration) {
		return func(*gin.Context) (*rate.Limiter, time.Duration) {
			return rate.NewLimiter(rate.Every(reqDelay), limit), time.Hour
		}
	}
	rateLimitAbortFn = func(c *gin.Context) {
		app := ginApp.GetApp(c)
		app.RateLimitError()
	}
)

type KeyGenerateFn func(*gin.Context) string
type AbortFn func(*gin.Context)

func KeyRateLimit(key KeyGenerateFn, createLimiter func(*gin.Context) (*rate.Limiter,
	time.Duration), abortFn AbortFn) gin.HandlerFunc {
	return func(c *gin.Context) {
		k := key(c)
		limiterInfo, ok := limiterSet.Get(k)
		if !ok {
			var expire time.Duration
			limiterInfo, expire = createLimiter(c)
			limiterSet.Set(k, limiterInfo, expire)
		}
		if ok = limiterInfo.(*rate.Limiter).Allow(); !ok {
			abortFn(c)
			return
		}
		c.Next()
	}
}
func RateLimit(lmt *limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		app := ginApp.GetApp(c)
		httpError := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if httpError != nil {
			app.RateLimitError()
			return
		}
		c.Next()
	}
}

// 每个ip每多久时间可以请求一次
func ClientIPRateLimit(reqDelay time.Duration, limitTimesParam ...int) gin.HandlerFunc {
	limit := 1
	if len(limitTimesParam) > 0 {
		limit = limitTimesParam[0]
	}
	return KeyRateLimit(getReqClientIP, getLimitFn(reqDelay, limit), rateLimitAbortFn)
}
