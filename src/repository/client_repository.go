package repository

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"go-bank/src/model"
	"log"
)

func CreateClient(db *gorm.DB, client model.Client) (model.Client, error) {

	tx := db.Begin()
	if tx.Error != nil {
		log.Fatalf("Erro ao iniciar a transação: %v", tx.Error)
	}

	if err := tx.Save(&client).Error; err != nil {
		tx.Rollback()
		log.Fatalf("Erro ao criar usuario: %v", err)
		return model.Client{}, nil
	}

	if err := tx.Commit().Error; err != nil {
		log.Fatalf("Erro ao commitar a transação: %v", err)
	}
	return client, nil
}

func ListClient(db *gorm.DB) (*[]model.Client, error) {
	var clients []model.Client
	db.Order("age desc").Order("name").Find(&clients)
	fmt.Printf("Clients not found: %+v\n", clients)
	return &clients, nil
}
