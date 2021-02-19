package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/real-web-world/go-web-api/global"
	ginApp "github.com/real-web-world/go-web-api/pkg/gin"
)

func Logger(c *gin.Context) {
	c.Next()
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
