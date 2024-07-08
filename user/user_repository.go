package user

import (
	"github.com/jinzhu/gorm"
	"log"
)

func Create(db *gorm.DB, user User) (User, error) {

	tx := db.Begin()
	if tx.Error != nil {
		log.Fatalf("Erro ao iniciar a transação: %v", tx.Error)
	}

	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		log.Fatalf("Erro ao criar usuario: %v", err)
		return User{}, nil
	}

	if err := tx.Commit().Error; err != nil {
		log.Fatalf("Erro ao commitar a transação: %v", err)
	}
	return user, nil
}
