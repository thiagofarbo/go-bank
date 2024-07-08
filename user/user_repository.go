package user

import (
	"github.com/jinzhu/gorm"
	"log"
)

func Create(db *gorm.DB, user User) (User, error) {

	tx := db.Begin()
	if tx.Error != nil {
		log.Fatalf("Error to start transaction: %v", tx.Error)
	}

	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		log.Fatalf("Error to create user: %v", err)
		return User{}, nil
	}

	if err := tx.Commit().Error; err != nil {
		log.Fatalf("Error to commit transaction: %v", err)
	}
	return user, nil
}
