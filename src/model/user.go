package model

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Username  string    `gorm:"size:64"`
	Password  string    `gorm:"size:255"`
	Email     string    `gorm:"size:100;not null;unique"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
