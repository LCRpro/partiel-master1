package controllers

import (
	"encoding/json"
	"strconv"
	"partiel-master1/config"
	"partiel-master1/models"

	"github.com/gin-gonic/gin"
)

func getUserEmailAndRoles(c *gin.Context) (string, []string, bool) {
	emailAny, exists := c.Get("email")
	if !exists {
		return "", nil, false
	}
	email := emailAny.(string)
	var user models.User
	result := config.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return "", nil, false
	}
	var roles []string
	json.Unmarshal([]byte(user.Roles), &roles)
	return email, roles, true
}

func isSpeaker(roles []string) bool {
	for _, r := range roles {
		if r == "conferencier" {
			return true
		}
	}
	return false
}

func CreateConference(c *gin.Context) {
	email, roles, ok := getUserEmailAndRoles(c)
	if !ok || !isSpeaker(roles) {
		c.JSON(403, gin.H{"error": "Seuls les conférenciers peuvent créer une conférence"})
		return
	}

	var req models.Conference
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Données invalides"})
		return
	}

	if req.StartTime == "" || req.EndTime == "" || req.StartTime == req.EndTime {
		c.JSON(400, gin.H{"error": "Durée invalide"})
		return
	}

	var count int64
	config.DB.Model(&models.Conference{}).
		Where("room = ? AND ((start_time < ? AND end_time > ?) OR (start_time < ? AND end_time > ?))",
			req.Room, req.EndTime, req.StartTime, req.StartTime, req.EndTime).
		Count(&count)
	if count > 0 {
		c.JSON(400, gin.H{"error": "Salle déjà occupée sur ce créneau"})
		return
	}

	day := req.StartTime[:10]
	var dayCount int64
	config.DB.Model(&models.Conference{}).
		Where("room = ? AND DATE(start_time) = ?", req.Room, day).
		Count(&dayCount)
	if dayCount >= 10 {
		c.JSON(400, gin.H{"error": "Limite de 10 conférences par salle et par jour atteinte"})
		return
	}

	var user models.User
	config.DB.Where("email = ?", email).First(&user)
	req.OrganizerID = user.ID
	req.SpeakerName = user.SpeakerName

	config.DB.Create(&req)
	c.JSON(201, req)
}

func ListConferences(c *gin.Context) {
	var confs []models.Conference
	config.DB.Find(&confs)

	type ConfWithCount struct {
		models.Conference
		ParticipantCount int `json:"ParticipantCount"`
	}

	var result []ConfWithCount
	for _, conf := range confs {
		var count int64
		config.DB.Model(&models.UserConference{}).Where("conference_id = ?", conf.ID).Count(&count)
		result = append(result, ConfWithCount{
			Conference:       conf,
			ParticipantCount: int(count),
		})
	}

	c.JSON(200, result)
}

func JoinConference(c *gin.Context) {
	email, _, ok := getUserEmailAndRoles(c)
	if !ok {
		c.JSON(401, gin.H{"error": "Non authentifié"})
		return
	}

	var user models.User
	config.DB.Where("email = ?", email).First(&user)
	confID := c.Param("id")

	var existing models.UserConference
	config.DB.Where("user_id = ? AND conference_id = ?", user.ID, confID).First(&existing)
	if existing.ID != 0 {
		c.JSON(400, gin.H{"error": "Déjà inscrit"})
		return
	}

	var count int64
	config.DB.Model(&models.UserConference{}).Where("conference_id = ?", confID).Count(&count)
	if count >= 20 {
		c.JSON(400, gin.H{"error": "Conférence complète"})
		return
	}

	insc := models.UserConference{
		UserID:       user.ID,
		ConferenceID: uint(ParseUint(confID)),
	}
	config.DB.Create(&insc)
	c.JSON(200, gin.H{"message": "Inscrit à la conférence"})
}

func ListMyConferences(c *gin.Context) {
	email, _, ok := getUserEmailAndRoles(c)
	if !ok {
		c.JSON(401, gin.H{"error": "Non authentifié"})
		return
	}
	var user models.User
	config.DB.Where("email = ?", email).First(&user)

	var inscs []models.UserConference
	config.DB.Where("user_id = ?", user.ID).Find(&inscs)

	var confIDs []uint
	for _, i := range inscs {
		confIDs = append(confIDs, i.ConferenceID)
	}

	var confs []models.Conference
	config.DB.Where("id IN ?", confIDs).Find(&confs)
	c.JSON(200, confs)
}

func ParseUint(s string) uint {
	u, _ := strconv.ParseUint(s, 10, 64)
	return uint(u)
}
