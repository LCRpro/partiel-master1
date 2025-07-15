package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"partiel-master1/config"
	"partiel-master1/models"
)

func insertTestUser(email string, roles string) models.User {
	user := models.User{
		GoogleID: "gid",
		Email:    email,
		Name:     "TestUser",
		Roles:    roles,
	}
	config.DB.Create(&user)
	return user
}

func Test_GetMe_NonAuth(t *testing.T) {
	cleanDB()
	router := gin.New()
	router.GET("/me", GetMe)
	req, _ := http.NewRequest("GET", "/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)
}

func Test_GetMe_NotFound(t *testing.T) {
	cleanDB()
	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set("email", "notfound@mail.com"); c.Next() })
	router.GET("/me", GetMe)
	req, _ := http.NewRequest("GET", "/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}

func Test_GetMe_Ok(t *testing.T) {
	cleanDB()
	user := insertTestUser("me@mail.com", `["user"]`)
	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set("email", user.Email); c.Next() })
	router.GET("/me", GetMe)
	req, _ := http.NewRequest("GET", "/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "me@mail.com")
}

func Test_UpdateRole_NonAuth(t *testing.T) {
	cleanDB()
	router := gin.New()
	router.PATCH("/role", UpdateRole)
	req, _ := http.NewRequest("PATCH", "/role", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)
}

func Test_UpdateRole_UserNotFound(t *testing.T) {
	cleanDB()
	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set("email", "notfound@mail.com"); c.Next() })
	router.PATCH("/role", UpdateRole)
	payload := map[string]interface{}{"role": "admin"}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PATCH", "/role", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}

func Test_UpdateRole_BadPayload(t *testing.T) {
	cleanDB()
	user := insertTestUser("test@mail.com", `["user"]`)
	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set("email", user.Email); c.Next() })
	router.PATCH("/role", UpdateRole)
	payload := map[string]interface{}{"role": "badrole"}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PATCH", "/role", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func Test_UpdateRole_AddAdmin(t *testing.T) {
	cleanDB()
	user := insertTestUser("admin@mail.com", `["user"]`)
	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set("email", user.Email); c.Next() })
	router.PATCH("/role", UpdateRole)
	payload := map[string]interface{}{"role": "admin"}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PATCH", "/role", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "admin")
}

func Test_UpdateRole_AddConferencierWithSpeakerName(t *testing.T) {
	cleanDB()
	user := insertTestUser("speaker@mail.com", `["user"]`)
	router := gin.New()
	router.Use(func(c *gin.Context) { c.Set("email", user.Email); c.Next() })
	router.PATCH("/role", UpdateRole)
	payload := map[string]interface{}{
		"role": "conferencier", "speakerName": "SpeakerName",
	}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PATCH", "/role", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "conferencier")
	assert.Contains(t, w.Body.String(), "SpeakerName")
}

func countRoleInJSON(jsonStr, role string) int {
	type resp struct {
		User struct {
			Roles string
		}
	}
	var r resp
	_ = json.Unmarshal([]byte(jsonStr), &r)
	var arr []string
	json.Unmarshal([]byte(r.User.Roles), &arr)
	n := 0
	for _, x := range arr {
		if x == role {
			n++
		}
	}
	return n
}
