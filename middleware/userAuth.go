package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/real-web-world/go-api/pkg/gin"
)

func AdminAuth(c *gin.Context) {
	app := ginApp.GetApp(c)
	if !app.IsAdmin {
		app.NoLogin()
		return
	}
	c.Next()
}
func CommonUserAuth(c *gin.Context) {
	app := ginApp.GetApp(c)
	if !app.IsLogin || app.AuthUser.Level != ginApp.LevelGeneral {
		app.NoLogin()
		return
	}
	c.Next()
}
func AuthorUserAuth(c *gin.Context) {
	app := ginApp.GetApp(c)
	if !app.IsLogin || app.AuthUser.Level != ginApp.LevelAuthor {
		app.NoLogin()
		return
	}
	c.Next()
}
func AdminOrAuthorAuth(c *gin.Context) {
	app := ginApp.GetApp(c)
	if !app.IsLogin || (!app.IsSuper && !app.IsAdmin && !app.IsAuthor) {
		app.NoLogin()
		return
	}
	c.Next()
}
