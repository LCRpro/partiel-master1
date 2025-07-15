package controllers

import (
	"partiel-master1/config"
	"partiel-master1/models"
	"encoding/json"
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



func UpdateRole(c *gin.Context) {
	emailAny, exists := c.Get("email")
	if !exists {
		c.JSON(401, gin.H{"error": "Non authentifié"})
		return
	}
	email := emailAny.(string)

	var req struct {
		NewRole string `json:"role"` 
	}
	if err := c.ShouldBindJSON(&req); err != nil || (req.NewRole != "admin" && req.NewRole != "organisateur") {
		c.JSON(400, gin.H{"error": "Role invalide"})
		return
	}

	var user models.User
	result := config.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		c.JSON(404, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	var roles []string
	if user.Roles != "" {
		json.Unmarshal([]byte(user.Roles), &roles)
	}
	found := false
	for _, r := range roles {
		if r == req.NewRole {
			found = true
			break
		}
	}
	if !found {
		roles = append(roles, req.NewRole)
	}

	rolesJSON, _ := json.Marshal(roles)
	user.Roles = string(rolesJSON)
	config.DB.Save(&user)

	c.JSON(200, user)
}

