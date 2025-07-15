package routes

import (
	"partiel-master1/controllers"
	"partiel-master1/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/auth/google/login", controllers.GoogleLogin)
	r.GET("/auth/google/callback", controllers.GoogleCallback)

	secured := r.Group("/")
	secured.Use(middlewares.JWTAuth())
	secured.GET("/me", controllers.GetMe)
}
