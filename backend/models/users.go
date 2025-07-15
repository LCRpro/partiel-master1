package models


type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	GoogleID string `gorm:"uniqueIndex"` 
	Email    string `gorm:"uniqueIndex"`
	Name     string
	Picture  string
}
