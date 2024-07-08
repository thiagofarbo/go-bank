package user

import (
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"size:100;not null"`
	Age       uint      `gorm:"not null"`
	Email     string    `gorm:"size:100;not null;unique"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func CreateUser(db *gorm.DB, user User) (User, error) {

	newUser := User{Name: user.Name, Age: user.Age, Email: user.Email}

	newUser, _ = Create(db, user)

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
