package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"partiel-master1/config"
	"partiel-master1/models"
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
		NewRole     string `json:"role"`
		SpeakerName string `json:"speakerName"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || (req.NewRole != "admin" && req.NewRole != "organisateur" && req.NewRole != "conferencier") {
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
	if req.NewRole == "conferencier" && req.SpeakerName != "" {
		user.SpeakerName = req.SpeakerName
	}
	rolesJSON, _ := json.Marshal(roles)
	user.Roles = string(rolesJSON)
	config.DB.Save(&user)

	jwtToken, _ := generateJWT(user.Email)

	c.JSON(200, gin.H{
		"user":  user,
		"token": jwtToken,
	})
}
