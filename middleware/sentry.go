package middleware

import (
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	sentryGin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"

	"github.com/real-web-world/go-web-api/conf"
)

func GenerateSentryMiddleware(appConf *conf.AppConf) gin.HandlerFunc {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn: appConf.Sentry.Dsn,
	}); err != nil {
		log.Fatalln("Sentry initialization failed: ", err)
	}
	return sentryGin.New(sentryGin.Options{
		Repanic: true,
		Timeout: 3 * time.Second,
	})
}
