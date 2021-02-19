package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/real-web-world/go-web-api/global"
	"github.com/real-web-world/go-web-api/models"
	"github.com/real-web-world/go-web-api/pkg/gin"
	"github.com/real-web-world/go-web-api/services/cache"
)

func Auth(c *gin.Context) {
	app := ginApp.GetApp(c)
	authUser := &ginApp.AuthUser{}
	token := app.GetToken()
	if token != "" {
		if uid, err := cache.GetUIDByToken(token); err == nil {
			user := models.NewCtxUser(c)
			if err := models.Get(user, uid); err == nil {
				authUser.IsLogin = true
				authUser.Level = user.Level
				authUser.Name = user.Name
				authUser.User = user
				authUser.UID = user.ID
				app.SetCtxAuthUser(authUser)
				app.SetUser(authUser)
				go func() { _ = cache.UpdateTokenExpire(token, global.Conf.Token.Expire*60) }()
			}
		}
	}
	c.Next()
}
