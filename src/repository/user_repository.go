package repository

import (
	"github.com/jinzhu/gorm"
	"go-bank/src/model"
	"log"
)

func CreateUser(db *gorm.DB, user model.User) (model.User, error) {

	tx := db.Begin()
	if tx.Error != nil {
		log.Fatalf("Error to start transaction: %v", tx.Error)
	}

	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		log.Fatalf("Error to create user: %v", err)
		return user, nil
	}

	if err := tx.Commit().Error; err != nil {
		log.Fatalf("Error to commit transaction: %v", err)
	}
	return user, nil
}

//func GetUserById(db *gorm.DB, id uint) (model.User, error) {
//	var user model.User
//	result := db.First(&user, id)
//	if result.Error != nil {
//		return model.User{}, nil
//	}
//	return user, nil
//}
