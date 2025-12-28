package models

import "time"


type URL struct {
	ShortCode string    `json:"short_code" gorm:"primaryKey;column:short_code"`     
	LongURL   string    `json:"long_url" gorm:"column:long_url;not null"`           
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"` 
	Clicks    int64     `json:"clicks" gorm:"column:clicks;default:0"`              
}


func (URL) TableName() string {
	return "urls"
}
