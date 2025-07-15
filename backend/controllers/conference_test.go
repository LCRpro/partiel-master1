package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"


	"partiel-master1/models"
	"partiel-master1/config"
)




func insertSpeaker() models.User {
	user := models.User{
		GoogleID:    "fake-google-id-speaker",
		Email:       "speaker@mail.com",
		Name:        "Speaker",
		SpeakerName: "JohnConf",
		Roles:       `["user", "conferencier"]`,
	}
	config.DB.Create(&user)
	return user
}

func insertUser() models.User {
	user := models.User{
		GoogleID:    "fake-google-id-user",
		Email:       "user@mail.com",
		Name:        "User",
		Roles:       `["user"]`,
	}
	config.DB.Create(&user)
	return user
}




func Test_isSpeaker(t *testing.T) {
	assert.True(t, isSpeaker([]string{"user", "conferencier"}))
	assert.False(t, isSpeaker([]string{"user", "admin"}))
	assert.False(t, isSpeaker([]string{}))
}

func Test_CreateConference_RefuseSiNonAuth(t *testing.T) {
	cleanDB()
	router := gin.New()
	router.POST("/conferences", CreateConference)
	payload := map[string]interface{}{
		"Title": "BadConf", "Description": "desc", "Room": 1,
		"StartTime": "2025-07-18 10:00", "EndTime": "2025-07-18 11:00",
	}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/conferences", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 403, w.Code)
}

func Test_CreateConference_RefuseSiPasSpeaker(t *testing.T) {
	cleanDB()
	user := insertUser()
	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set("email", user.Email); c.Next() })
	router.POST("/conferences", CreateConference)
	payload := map[string]interface{}{
		"Title": "BadConf", "Description": "desc", "Room": 1,
		"StartTime": "2025-07-18 10:00", "EndTime": "2025-07-18 11:00",
	}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/conferences", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 403, w.Code)
}

func Test_CreateConference_DonneesInvalides(t *testing.T) {
	cleanDB()
	speaker := insertSpeaker()
	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set("email", speaker.Email); c.Next() })
	router.POST("/conferences", CreateConference)
	req, _ := http.NewRequest("POST", "/conferences", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func Test_CreateConference_DureeInvalide(t *testing.T) {
	cleanDB()
	speaker := insertSpeaker()
	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set("email", speaker.Email); c.Next() })
	router.POST("/conferences", CreateConference)
	payload := map[string]interface{}{
		"Title": "Bad", "Description": "desc", "Room": 1,
		"StartTime": "", "EndTime": "",
	}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/conferences", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func Test_CreateConference_Ok(t *testing.T) {
	cleanDB()
	speaker := insertSpeaker()
	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set("email", speaker.Email); c.Next() })
	router.POST("/conferences", CreateConference)
	payload := map[string]interface{}{
		"Title": "GoodConf", "Description": "desc", "Room": 1,
		"StartTime": "2025-07-18 10:00", "EndTime": "2025-07-18 11:00",
	}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/conferences", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
	assert.Contains(t, w.Body.String(), "GoodConf")
}

func Test_ListConferences_Empty(t *testing.T) {
	cleanDB()
	router := gin.New()
	router.GET("/confs", ListConferences)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/confs", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	var res []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Len(t, res, 0)
}

func Test_ListConferences_NotEmpty(t *testing.T) {
	cleanDB()
	speaker := insertSpeaker()
	conf := models.Conference{
		Title: "T1", Description: "desc", Room: 1,
		StartTime: "2025-07-18 10:00", EndTime: "2025-07-18 11:00",
		OrganizerID: speaker.ID, SpeakerName: speaker.SpeakerName,
	}
	config.DB.Create(&conf)
	router := gin.New()
	router.GET("/confs", ListConferences)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/confs", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	var res []map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &res)
	assert.Len(t, res, 1)
}

func Test_JoinConference_Ok(t *testing.T) {
	cleanDB()
	speaker := insertSpeaker()
	user := insertUser()
	conf := models.Conference{
		Title: "ConfToJoin", Room: 1, StartTime: "2025-07-18 10:00", EndTime: "2025-07-18 11:00",
		OrganizerID: speaker.ID, SpeakerName: speaker.SpeakerName,
	}
	config.DB.Create(&conf)

	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set("email", user.Email); c.Next() })
	router.POST("/conferences/:id/join", JoinConference)
	req, _ := http.NewRequest("POST", "/conferences/"+strconv.Itoa(int(conf.ID))+"/join", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Inscrit")
}

func Test_JoinConference_DejaInscrit(t *testing.T) {
	cleanDB()
	speaker := insertSpeaker()
	user := insertUser()
	conf := models.Conference{
		Title: "ConfToJoin", Room: 1, StartTime: "2025-07-18 10:00", EndTime: "2025-07-18 11:00",
		OrganizerID: speaker.ID, SpeakerName: speaker.SpeakerName,
	}
	config.DB.Create(&conf)
	insc := models.UserConference{UserID: user.ID, ConferenceID: conf.ID}
	config.DB.Create(&insc)

	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set("email", user.Email); c.Next() })
	router.POST("/conferences/:id/join", JoinConference)
	req, _ := http.NewRequest("POST", "/conferences/"+strconv.Itoa(int(conf.ID))+"/join", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "Déjà inscrit")
}

func Test_JoinConference_NonAuth(t *testing.T) {
	cleanDB()
	router := gin.New()
	router.POST("/conferences/:id/join", JoinConference)
	req, _ := http.NewRequest("POST", "/conferences/1/join", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)
}

func Test_ListMyConferences_NonAuth(t *testing.T) {
	cleanDB()
	router := gin.New()
	router.GET("/myconfs", ListMyConferences)
	req, _ := http.NewRequest("GET", "/myconfs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)
}

func Test_ListMyConferences_Ok(t *testing.T) {
	cleanDB()
	user := insertUser()
	speaker := insertSpeaker()
	conf := models.Conference{
		Title: "Conf2", Room: 2, StartTime: "2025-07-18 10:00", EndTime: "2025-07-18 11:00",
		OrganizerID: speaker.ID, SpeakerName: speaker.SpeakerName,
	}
	config.DB.Create(&conf)
	insc := models.UserConference{UserID: user.ID, ConferenceID: conf.ID}
	config.DB.Create(&insc)

	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set("email", user.Email); c.Next() })
	router.GET("/myconfs", ListMyConferences)
	req, _ := http.NewRequest("GET", "/myconfs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	var confs []models.Conference
	json.Unmarshal(w.Body.Bytes(), &confs)
	assert.Len(t, confs, 1)
}

func Test_UpdateConference_NonAuth(t *testing.T) {
	cleanDB()
	router := gin.New()
	router.PATCH("/conferences/:id", UpdateConference)
	req, _ := http.NewRequest("PATCH", "/conferences/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)
}

func Test_UpdateConference_PasOrga(t *testing.T) {
	cleanDB()
	speaker := insertSpeaker()
	user := insertUser()
	conf := models.Conference{
		Title: "Conf2", Room: 2, StartTime: "2025-07-18 10:00", EndTime: "2025-07-18 11:00",
		OrganizerID: speaker.ID, SpeakerName: speaker.SpeakerName,
	}
	config.DB.Create(&conf)
	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set("email", user.Email); c.Next() })
	router.PATCH("/conferences/:id", UpdateConference)
	req, _ := http.NewRequest("PATCH", "/conferences/"+strconv.Itoa(int(conf.ID)), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 403, w.Code)
}

func Test_UpdateConference_OK(t *testing.T) {
	cleanDB()
	speaker := insertSpeaker()
	conf := models.Conference{
		Title: "ConfOld", Room: 1, StartTime: "2025-07-18 10:00", EndTime: "2025-07-18 11:00",
		OrganizerID: speaker.ID, SpeakerName: speaker.SpeakerName,
	}
	config.DB.Create(&conf)
	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set("email", speaker.Email); c.Next() })
	router.PATCH("/conferences/:id", UpdateConference)
	payload := map[string]interface{}{
		"Title": "ConfNew",
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PATCH", "/conferences/"+strconv.Itoa(int(conf.ID)), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "ConfNew")
}

func Test_DeleteConference_OK(t *testing.T) {
	cleanDB()
	speaker := insertSpeaker()
	conf := models.Conference{
		Title: "ConfDEL", Room: 1, StartTime: "2025-07-18 10:00", EndTime: "2025-07-18 11:00",
		OrganizerID: speaker.ID, SpeakerName: speaker.SpeakerName,
	}
	config.DB.Create(&conf)
	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set("email", speaker.Email); c.Next() })
	router.DELETE("/conferences/:id", DeleteConference)
	req, _ := http.NewRequest("DELETE", "/conferences/"+strconv.Itoa(int(conf.ID)), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "supprimée")
}

func Test_DeleteConference_NonAuth(t *testing.T) {
	cleanDB()
	router := gin.New()
	router.DELETE("/conferences/1", DeleteConference)
	req, _ := http.NewRequest("DELETE", "/conferences/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)
}

func Test_DeleteConference_PasOrga(t *testing.T) {
	cleanDB()
	speaker := insertSpeaker()
	user := insertUser()
	conf := models.Conference{
		Title: "ConfDEL", Room: 1, StartTime: "2025-07-18 10:00", EndTime: "2025-07-18 11:00",
		OrganizerID: speaker.ID, SpeakerName: speaker.SpeakerName,
	}
	config.DB.Create(&conf)
	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set("email", user.Email); c.Next() })
	router.DELETE("/conferences/"+strconv.Itoa(int(conf.ID)), DeleteConference)
	req, _ := http.NewRequest("DELETE", "/conferences/"+strconv.Itoa(int(conf.ID)), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 403, w.Code)
}

func Test_ListConferenceParticipants_Ok(t *testing.T) {
	cleanDB()
	speaker := insertSpeaker()
	user := insertUser()
	conf := models.Conference{
		Title: "Conf", Room: 1, StartTime: "2025-07-18 10:00", EndTime: "2025-07-18 11:00",
		OrganizerID: speaker.ID, SpeakerName: speaker.SpeakerName,
	}
	config.DB.Create(&conf)
	insc := models.UserConference{UserID: user.ID, ConferenceID: conf.ID}
	config.DB.Create(&insc)
	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set("email", speaker.Email); c.Next() })
	router.GET("/conferences/:id/participants", ListConferenceParticipants)
	req, _ := http.NewRequest("GET", "/conferences/"+strconv.Itoa(int(conf.ID))+"/participants", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), user.Email)
}

func Test_ListConferenceParticipants_NonAuth(t *testing.T) {
	cleanDB()
	speaker := insertSpeaker()
	conf := models.Conference{
		Title: "Conf", Room: 1, StartTime: "2025-07-18 10:00", EndTime: "2025-07-18 11:00",
		OrganizerID: speaker.ID, SpeakerName: speaker.SpeakerName,
	}
	config.DB.Create(&conf)
	router := gin.New()
	router.GET("/conferences/:id/participants", ListConferenceParticipants)
	req, _ := http.NewRequest("GET", "/conferences/"+strconv.Itoa(int(conf.ID))+"/participants", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)
}

func Test_ListConferenceParticipants_PasOrga(t *testing.T) {
	cleanDB()
	speaker := insertSpeaker()
	user := insertUser()
	conf := models.Conference{
		Title: "Conf", Room: 1, StartTime: "2025-07-18 10:00", EndTime: "2025-07-18 11:00",
		OrganizerID: speaker.ID, SpeakerName: speaker.SpeakerName,
	}
	config.DB.Create(&conf)
	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set("email", user.Email); c.Next() })
	router.GET("/conferences/:id/participants", ListConferenceParticipants)
	req, _ := http.NewRequest("GET", "/conferences/"+strconv.Itoa(int(conf.ID))+"/participants", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 403, w.Code)
}

func Test_ListConferenceParticipants_NotFound(t *testing.T) {
	cleanDB()
	speaker := insertSpeaker()
	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set("email", speaker.Email); c.Next() })
	router.GET("/conferences/999/participants", ListConferenceParticipants)
	req, _ := http.NewRequest("GET", "/conferences/999/participants", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}
