package routes

import (
	"github.com/LCRpro/partiel-master1/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/auth/google/login", controllers.GoogleLogin)
	router.GET("/auth/google/callback", controllers.GoogleCallback)
}
