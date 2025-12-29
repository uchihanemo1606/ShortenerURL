package models

import "time"

type User struct {
	ID			string      `json:"id" gorm:"primaryKey;autoIncrement"`
	Email   	string	`json:"email" gorm:"uniqueIndex;not null"`
	Password 	string    `json:"password" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}