package controllers

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGenerateJWT(t *testing.T) {
	os.Setenv("JWT_SECRET", "TEST_SECRET")
	token, err := generateJWT("test@example.com")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGoogleResponseUnmarshal(t *testing.T) {
	data := `{
        "id": "123456",
        "email": "mail@mail.com",
        "name": "Jane",
        "picture": "http://img"
    }`
	var gr GoogleResponse
	err := json.Unmarshal([]byte(data), &gr)
	assert.NoError(t, err)
	assert.Equal(t, "Jane", gr.Name)
	assert.Equal(t, "mail@mail.com", gr.Email)
}
