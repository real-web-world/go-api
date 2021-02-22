package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/real-web-world/go-api/pkg/gin"
)

func NotSaveResp(c *gin.Context) {
	app := ginApp.GetApp(c)
	app.SetNotSaveResp()
	c.Next()
}
