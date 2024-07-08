package grossIncome

import (
	"github.com/jinzhu/gorm"
	"go-bank/account"
	"log"
	"time"
)

type GrossIncome struct {
	ID        uint            `gorm:"primaryKey;autoIncrement"`
	AccountID uint            `gorm:"not null"`
	Account   account.Account `gorm:"foreignkey:AccountID;association_foreignkey:ID"`
	Amount    float64         `gorm:"type:decimal(15,2);not null"`
	CreatedAt time.Time       `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time       `gorm:"default:CURRENT_TIMESTAMP"`
}

func Create(db *gorm.DB, income GrossIncome) (GrossIncome, error) {

	tx := db.Begin()
	if tx.Error != nil {
		log.Fatalf("Error to start transaction: %v", tx.Error)
	}

	if err := tx.Save(&income).Error; err != nil {
		tx.Rollback()
		log.Fatalf("Error to create Income: %v", err)
		return GrossIncome{}, nil
	}

	if err := tx.Commit().Error; err != nil {
		log.Fatalf("Error to commit transaction: %v", err)
	}
	return income, nil
}
