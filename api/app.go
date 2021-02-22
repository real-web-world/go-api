package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/real-web-world/go-api/pkg/gin"
)

func Index(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if c.Query("buff") == "buff" {
		app := ginApp.GetApp(c)
		app.String("Hello buff")
		return
	}
	c.File("./public/index.html")
}
func ServerInfo(c *gin.Context) {
	app := ginApp.GetApp(c)
	app.Data(ginApp.ServerInfo{Timestamp: time.Now().Unix()})
}
