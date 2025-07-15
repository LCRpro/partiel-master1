package models

import "gorm.io/gorm"

var DB *gorm.DB

func SetTestDB(db *gorm.DB) {
	DB = db
}
