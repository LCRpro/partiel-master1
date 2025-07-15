package models


type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	GoogleID string `gorm:"type:varchar(64);uniqueIndex"`
	Email    string `gorm:"type:varchar(255);uniqueIndex"`
	Name     string
	Picture  string
	Roles    string `gorm:"type:json"`
}

