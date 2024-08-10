package client

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Client struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"size:100;not null"`
	Age       uint      `gorm:"not null"`
	Email     string    `gorm:"size:100;not null;unique"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

//func CreateClient(db *gorm.DB, client Client) (Client, error) {
//
//	newUser := Client{Name: client.Name, Age: client.Age, Email: client.Email}
//
//	newUser, _ = Create(db, client)
//
//	return newUser, nil
//}

func GetClientById(db *gorm.DB, id uint) (Client, error) {
	var user Client
	result := db.First(&user, id)
	if result.Error != nil {
		return Client{}, nil
	}
	return user, nil
}
