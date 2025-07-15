package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/LCRpro/partiel-master1/config"
	"github.com/LCRpro/partiel-master1/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erreur de chargement du .env")
	}

	config.InitOAuth()

	router := gin.Default()
	routes.SetupRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
