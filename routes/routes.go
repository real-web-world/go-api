package routes

import "github.com/gin-gonic/gin"

func RouterSetup(r *gin.Engine) {
	initWebRoutes(r)
	initAPIRoutes(r)
}
