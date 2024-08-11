package client

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type Client struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"size:100;not null"`
	Age       string    `gorm:"not null"`
	Email     string    `gorm:"size:100;not null;unique"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func GetClientById(db *gorm.DB, id uint) (Client, error) {
	var user Client
	result := db.First(&user, id)
	if result.Error != nil {
		return Client{}, nil
	}
	return user, nil
}

func ListClient(db *gorm.DB) (*[]Client, error) {
	var clients []Client
	db.Order("age desc").Order("name").Find(&clients)
	fmt.Printf("Clients not found: %+v\n", clients)
	return &clients, nil
}
