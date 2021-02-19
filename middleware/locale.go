package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/real-web-world/go-web-api/conf"
	"github.com/real-web-world/go-web-api/global"
	"github.com/real-web-world/go-web-api/pkg/gin"
)

func Locale(appConf *conf.AppConf) gin.HandlerFunc {
	return func(c *gin.Context) {
		app := ginApp.GetApp(c)
		locale := app.GetLocale()
		if locale == "" {
			locale = appConf.Lang.DefaultLang
		}
		translator := global.GetTrans(locale)
		app.SetTranslator(translator)
		c.Next()
	}
}
