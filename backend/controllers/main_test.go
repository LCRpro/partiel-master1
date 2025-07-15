package controllers

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"partiel-master1/config"
	"partiel-master1/models"
	"testing"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("Impossible d'ouvrir la DB de test: " + err.Error())
	}
	_ = db.AutoMigrate(&models.User{}, &models.Conference{}, &models.UserConference{})
	return db
}

func TestMain(m *testing.M) {
	db := setupTestDB()
	models.DB = db
	config.DB = db
	os.Exit(m.Run())
}

func cleanDB() {
	config.DB.Exec("DELETE FROM user_conferences")
	config.DB.Exec("DELETE FROM conferences")
	config.DB.Exec("DELETE FROM users")
}
