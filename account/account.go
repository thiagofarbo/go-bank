// No pacote account

package account

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"go-bank/user"
	"time"
)

type Account struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Branch    string    `gorm:"size:20;not null;unique"`
	UserID    uint      `gorm:"not null"`
	Number    string    `gorm:"size:20;not null;unique"`
	Balance   float64   `gorm:"type:numeric(10,2);not null;default:0.00"`
	Status    string    `gorm:"size:25;not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func CreateAccount(db *gorm.DB, account Account, user user.User) (Account, error) {
	account.UserID = user.ID
	result := db.Create(&account)
	if result.Error != nil {
		return Account{}, nil
	}
	return account, nil
}

func Deposit(db *gorm.DB, branch string, accountNumber string, amount float64) (*Account, error) {
	var account Account
	valid, err := IsValid(db, branch, accountNumber)
	if valid == false {
		return nil, err
	}

	db.Where("branch = ? AND number = ?", branch, accountNumber).Find(&account)
	if account.ID == 0 {
		fmt.Printf("Account not found : %v\n", branch)
		return &account, nil
	}
	account.Balance += amount
	db.Save(&account)
	return &account, nil
}

func Withdraw(db *gorm.DB, branch string, accountNumber string, amount float64) (*Account, error) {
	var account Account
	valid, err := IsValid(db, branch, accountNumber)
	if valid == false {
		return &Account{}, err
	}
	db.Where("branch = ? AND number = ?", branch, accountNumber).Find(&account)
	if account.ID == 0 {
		fmt.Printf("Account not found : %v\n", branch)
		return &account, nil
	}
	isAble, _ := IsAbleToWithdraw(account.Balance, amount)
	if isAble == false {
		return &Account{}, nil
	}
	account.Balance -= amount
	db.Save(&account)
	return &account, nil
}

func IsAbleToWithdraw(balance float64, withdrawAmount float64) (bool, error) {
	if balance < withdrawAmount {
		fmt.Printf("Unable to process withdrawal: insufficient funds : %v\n", balance)
		return false, nil
	}
	return true, nil
}

func IsValid(db *gorm.DB, branch string, accountNumber string) (bool, error) {
	var account Account
	db.Where("branch = ? AND number = ?", branch, accountNumber).Find(&account)
	if account.ID == 0 {
		fmt.Printf("Account not found : %v\n", branch)
		return false, nil
	}

	if account.Status != "active" {
		fmt.Printf("Account status error : %v\n", account.Status)
		return false, nil
	}
	return true, nil
}

func UpdateStatusAccount(db *gorm.DB, branch string, accountNumber string, status string) (*Account, error) {
	var account Account
	valid, err := IsValid(db, branch, accountNumber)
	if valid == false {
		return nil, err
	}
	db.Where("branch = ? AND number = ?", branch, accountNumber).Find(&account)
	if account.ID == 0 {
		fmt.Printf("Account not found : %v\n", branch)
		return &account, nil
	}
	account.Status = status
	db.Save(&account)
	return &account, nil
}

//func GetUserById(db *gorm.DB, id uint) (User, error) {
//	var user User
//	result := db.First(&user, id)
//	if result.Error != nil {
//		return User{}, nil
//	}
//	return user, nil
//}
