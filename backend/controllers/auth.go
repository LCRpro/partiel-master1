package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"partiel-master1/config"
	"partiel-master1/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type GoogleResponse struct {
	ID      string `json:"id"`       
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func GoogleLogin(c *gin.Context) {
	url := config.GoogleOauthConfig.AuthCodeURL("state")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	token, err := config.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token invalide"})
		return
	}

	client := config.GoogleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur Google user"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var googleData GoogleResponse
	err = json.Unmarshal(body, &googleData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur parse JSON Google"})
		return
	}

	user := models.User{
		GoogleID: googleData.ID,        
		Email:    googleData.Email,
		Name:     googleData.Name,
		Picture:  googleData.Picture,
	}

	var existing models.User
	result := config.DB.Where("email = ?", user.Email).First(&existing)

	if result.Error != nil {
		if err := config.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur insertion DB"})
			return
		}
	} else {
		user = existing 
	}

	jwtToken, _ := generateJWT(user.Email)

	userJson, _ := json.Marshal(user)
	redirect := fmt.Sprintf("http://localhost:5173/?token=%s&user=%s", jwtToken, url.QueryEscape(string(userJson)))
	c.Redirect(http.StatusTemporaryRedirect, redirect)
}

func generateJWT(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
