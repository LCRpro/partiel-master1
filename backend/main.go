package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"partiel-master1/config"
	"partiel-master1/models"
	"partiel-master1/routes"
)

func main() {
	godotenv.Load()
	config.InitOAuth()
	config.InitDB()

	models.Migrate(config.DB) 

	r := gin.Default()
	routes.SetupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
