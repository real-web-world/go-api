package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/real-web-world/go-web-api/pkg/gin"
	"github.com/real-web-world/go-web-api/services/cache"
)

func Cache(d time.Duration, keyParam ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Request.RequestURI
		if len(keyParam) == 1 {
			key = keyParam[0]
		}
		json := cache.GetAPICache(key)
		app := ginApp.GetApp(c)
		if json != nil {
			app.JSON(json)
			return
		}
		app.SetApiCacheKey(key)
		app.SetApiCacheExpire(&d)
		c.Next()
	}
}
