package models

import "time"

type User struct {
	ID           uint		`json:"user_id" example:"123`
	Username     string		`json:"username" example:"alice`
	PasswordHash string		
	CreatedAt    time.Time	`json:"created_at" example:"2025-11-14 10:23:45"`
}