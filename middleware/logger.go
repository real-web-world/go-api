package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/real-web-world/go-api/global"
	ginApp "github.com/real-web-world/go-api/pkg/gin"
)

func Logger(c *gin.Context) {
	c.Next()
	path := c.Request.URL.Path
	isDevApi := path == global.PrometheusApi ||
		strings.Index(path, global.DebugApiPrefix) == 0
	if isDevApi {
		return
	}
	app := ginApp.GetApp(c)
	sqls := app.GetSqls()
	reqID := app.GetReqID()
	var totalTime time.Duration
	for _, sql := range sqls {
		execTime, _ := time.ParseDuration(sql.ExecTime)
		totalTime += execTime
	}
	global.Logger.Infow("sql",
		"reqID", reqID,
		"sqlCount", len(sqls),
		"totalTime", totalTime.String(),
		"sqls", sqls,
	)
}
