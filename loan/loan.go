package loan

import (
	"go-bank/account"
	"go-bank/grossIncome"
	"time"
)

type Loan struct {
	ID           uint            `gorm:"primaryKey;autoIncrement"`
	AccountID    uint            `gorm:"not null"`
	Account      account.Account `gorm:"foreignkey:AccountID;association_foreignkey:ID"`
	Amount       float64         `gorm:"type:decimal(15,2);not null"`
	InterestRate float64         `gorm:"type:decimal(5,2);not null"`
	Term         int             `gorm:"not null"`
	Description  string          `gorm:"type:varchar(255)"`
	GrossIncomes []grossIncome.GrossIncome
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
