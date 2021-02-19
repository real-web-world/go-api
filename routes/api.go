package routes

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/real-web-world/go-web-api/api"
	mid "github.com/real-web-world/go-web-api/middleware"
	"github.com/real-web-world/go-web-api/models"
)

var (
	loginAuth                = mid.LoginAuth
	adminAuth                = mid.AdminAuth
	commonUserAuth           = mid.CommonUserAuth
	adminOrAuthorAuth        = mid.AdminOrAuthorAuth
	oneMinuteFiveTimeLimiter = mid.ClientIPRateLimit(time.Minute, 5)
	oneSecondLimiter         = mid.ClientIPRateLimit(time.Second, 1)
)

func initUserModule(r *gin.Engine) {
	user := r.Group("user")
	{
		user.POST("login",
			oneMinuteFiveTimeLimiter, oneSecondLimiter, api.Login)
		user.POST("logout", api.Logout)
		login := user.Group("/")
		login.Use(loginAuth)
		{
			login.POST("detail", api.UserDetail)
			login.POST("edit", api.EditUser)
			login.POST("modifyPwd", api.ModifyPwd)

		}
		admin := user.Group("/")
		admin.Use(adminAuth)
		{
			admin.POST("list", api.ListUser)
			admin.POST("add", api.AddUser)
			admin.POST("del", api.DelUser)
		}
	}
}
func initTagModule(r *gin.Engine) {
	m := r.Group("tag")
	{
		m.POST("detail", api.CommonDetail(&models.Tag{}))
		m.POST("list", api.CommonList(&models.Tag{}))
		admin := m.Group("/")
		admin.Use(adminAuth)
		{
			admin.POST("edit", api.CommonEdit(&models.Tag{}))
			admin.POST("add", api.CommonAdd(&models.Tag{}))
			admin.POST("del", api.CommonDel(&models.Tag{}))
		}
	}
}
func initCategoryModule(r *gin.Engine) {
	m := r.Group("category")
	{
		m.POST("detail", api.CommonDetail(&models.Category{}))
		m.POST("list", api.CommonList(&models.Category{}))
		admin := m.Group("/")
		admin.Use(adminAuth)
		{
			admin.POST("edit", api.CommonEdit(&models.Category{}))
			admin.POST("add", api.CommonAdd(&models.Category{}))
			admin.POST("del", api.CommonDel(&models.Category{}))
		}
	}
}
func initAPIRoutes(r *gin.Engine) {
	r.Any("test", api.TestHand)
	r.POST("/serverInfo", api.ServerInfo)
	initUserModule(r)
	initTagModule(r)
	initCategoryModule(r)
}
