package models


type UserConference struct {
	ID           uint `gorm:"primaryKey"`
	UserID       uint
	ConferenceID uint
}
