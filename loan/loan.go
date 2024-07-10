package loan

import (
	"github.com/jinzhu/gorm"
	"go-bank/account"
	"go-bank/grossIncome"
	"time"
)

type Loan struct {
	ID           uint            `gorm:"primaryKey;autoIncrement"`
	AccountID    uint            `gorm:"not null"`
	Account      account.Account `gorm:"foreignKey:AccountID"`
	Amount       float64         `gorm:"type:decimal(15,2);not null"`
	InterestRate float64         `gorm:"type:decimal(5,2);not null"`
	Term         int             `gorm:"not null"`
	Description  string          `gorm:"type:varchar(255)"`
	GrossIncomes []grossIncome.GrossIncome
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func GenerateLoan(db *gorm.DB, accountId uint, amount float64, interestRate float64, term int, description string) (Loan, error) {

	loan := Loan{AccountID: accountId, Amount: amount, InterestRate: interestRate, Term: term, Description: description + string(term)}

	result := db.Create(&loan)
	if result.Error != nil {
		return Loan{}, nil
	}
	return loan, nil
}
