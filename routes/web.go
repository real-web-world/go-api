package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/real-web-world/go-web-api/api"
	"github.com/real-web-world/go-web-api/global"
	mid "github.com/real-web-world/go-web-api/middleware"
)

func initWebRoutes(r *gin.Engine) {
	r.GET("/version", mid.NotSaveResp, api.ShowVersion)
	r.GET("/getVerifyCode", mid.NotSaveResp, api.GetCaptcha)
	r.GET(global.PrometheusApi, mid.NotSaveResp, gin.WrapH(promhttp.Handler()))
	r.NoRoute(mid.NotSaveResp, api.Index)
}
