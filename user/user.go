package user

import (
	"github.com/jinzhu/gorm"
	"go-bank/helper"
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

func CreateUser(db *gorm.DB, user User) (User, error) {

	hashPassword, err := helper.EncryptPassword(user.Password)
	if err != nil {
		return User{}, err
	}

	newUser := User{Username: user.Username, Password: string(hashPassword), Email: user.Email}

	newUser, _ = Create(db, newUser)

	return newUser, nil
}

func GetUserById(db *gorm.DB, id uint) (User, error) {
	var user User
	result := db.First(&user, id)
	if result.Error != nil {
		return User{}, nil
	}
	return user, nil
}
