package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/real-web-world/go-web-api/pkg/gin"
)

func LoginAuth(c *gin.Context) {
	app := ginApp.GetApp(c)
	if !app.IsLogin {
		app.NoLogin()
		return
	}
	c.Next()
}
