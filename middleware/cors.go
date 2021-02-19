package middleware

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/real-web-world/go-web-api/conf"
)

func Cors(appConf *conf.AppConf) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     strings.Split(appConf.Csrf.AllowOrigins, ","),
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT"},
		AllowHeaders:     []string{"content-type", "x-requested-with", "token", "locale"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
