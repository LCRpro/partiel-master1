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
	secured.POST("/me/role", controllers.UpdateRole)

	secured.POST("/conferences", controllers.CreateConference)
	secured.GET("/conferences", controllers.ListConferences)
	secured.POST("/conferences/:id/join", controllers.JoinConference)
	secured.GET("/me/conferences", controllers.ListMyConferences)
	secured.GET("/conferences/:id/participants", controllers.ListConferenceParticipants)
	secured.PATCH("/conferences/:id", controllers.UpdateConference)
	secured.DELETE("/conferences/:id", controllers.DeleteConference)

}
