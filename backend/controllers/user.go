package controllers

import (
	"partiel-master1/config"
	"partiel-master1/models"

	"github.com/gin-gonic/gin"
)

func GetMe(c *gin.Context) {
	emailAny, exists := c.Get("email")
	if !exists {
		c.JSON(401, gin.H{"error": "Non authentifié"})
		return
	}
	email := emailAny.(string)

	var user models.User
	result := config.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		c.JSON(404, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	c.JSON(200, user)
}
