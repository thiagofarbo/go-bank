package model

import (
	"github.com/jinzhu/gorm"
	"go-bank/src/helper"
	"time"
)

type Loan struct {
	ID           uint    `gorm:"primaryKey;autoIncrement"`
	AccountID    uint    `gorm:"not null"`
	Account      Account `gorm:"foreignKey:AccountID"`
	Amount       float64 `gorm:"type:decimal(15,2);not null"`
	InterestRate float64 `gorm:"type:decimal(5,2);not null"`
	Term         uint    `gorm:"not null"`
	Description  string  `gorm:"type:varchar(255)"`
	GrossIncomes []GrossIncome
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func GenerateLoan(db *gorm.DB, accountId uint, amount string, interestRate string, term string, description string) (Loan, error) {
	termInt, _ := helper.ToUint(term)
	amountFloat, _ := helper.ToFloat(amount)
	interestRateFloat, _ := helper.ToFloat(amount)
	loan := Loan{AccountID: accountId, Amount: amountFloat, InterestRate: interestRateFloat, Term: termInt, Description: description + term}

	result := db.Create(&loan)
	if result.Error != nil {
		return Loan{}, nil
	}
	return loan, nil
}
