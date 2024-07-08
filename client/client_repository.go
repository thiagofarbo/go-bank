package client

import (
	"github.com/jinzhu/gorm"
	"log"
)

func Create(db *gorm.DB, client Client) (Client, error) {

	tx := db.Begin()
	if tx.Error != nil {
		log.Fatalf("Erro ao iniciar a transação: %v", tx.Error)
	}

	if err := tx.Save(&client).Error; err != nil {
		tx.Rollback()
		log.Fatalf("Erro ao criar usuario: %v", err)
		return Client{}, nil
	}

	if err := tx.Commit().Error; err != nil {
		log.Fatalf("Erro ao commitar a transação: %v", err)
	}
	return client, nil
}
