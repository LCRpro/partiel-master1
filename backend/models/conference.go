package models


type Conference struct {
	ID           uint   `gorm:"primaryKey"`
	Title        string
	Description  string
	StartTime    string 
	EndTime      string
	Room         int    
	SpeakerName  string
	OrganizerID  uint   
}
