package models

import "time"


type URL struct {
	ShortCode string    `json:"short_code" gorm:"primaryKey;column:short_code;index"`     
	LongURL   string    `json:"long_url" gorm:"column:long_url;not null"`           
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"` 
	Clicks    int64     `json:"clicks" gorm:"column:clicks;default:0"`
	ExpiresAt time.Time `json:"expires_at" gorm:"column:expires_at"`
	UserID string       `json:"user_id" gorm:"column:user_id;not null;index"`
}


func (URL) TableName() string {
	return "urls"
}
