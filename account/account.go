// No pacote account

package account

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"go-bank/client"
	"go-bank/helper"
	"log"
	"time"
)

type Account struct {
	ID        uint          `gorm:"primaryKey;autoIncrement"`
	Branch    string        `gorm:"size:20;not null;unique"`
	ClientID  uint          `gorm:"not null"`
	Client    client.Client `gorm:"foreignKey:ClientID"`
	Number    string        `gorm:"size:20;not null;unique"`
	Balance   float64       `gorm:"type:numeric(10,2);not null;default:0.00"`
	Status    string        `gorm:"size:25;not null"`
	CreatedAt time.Time     `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time     `gorm:"default:CURRENT_TIMESTAMP"`
}

type Transaction struct {
	ID              uint      `gorm:"primaryKey;autoIncrement"`
	AccountID       uint      `gorm:"not null"`
	Amount          float64   `gorm:"type:numeric(10,2);not null;default:0.00"`
	TransactionType string    `gorm:"size:20;not null"`
	Description     string    `gorm:"size:20;not null"`
	Account         Account   `gorm:"foreignKey:AccountID"`
	CreatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func ListAccount(db *gorm.DB) (*[]Account, error) {
	var accounts []Account
	db.Order("branch desc").Order("number").Find(&accounts)
	fmt.Printf("Account not found: %+v\n", accounts)
	return &accounts, nil
}

func CreateAccount(db *gorm.DB, account Account, clientId uint) (Account, error) {
	account.ClientID = clientId
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

	_, err = CreateTransaction(db, amount, "Deposit", account)
	if err != nil {
		return nil, err
	}
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
		fmt.Printf("Unable to process withdraw: insufficient funds : %v\n", amount)
		return &Account{}, nil
	}
	account.Balance -= amount
	db.Save(&account)

	_, err = CreateTransaction(db, amount, "Withdraw", account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func IsAbleToWithdraw(balance float64, withdrawAmount float64) (bool, error) {
	if balance <= withdrawAmount {
		fmt.Printf("Unable to process withdrawal: insufficient funds : %v\n", balance)
		return false, nil
	}
	return true, nil
}

func IsAbleToTransfer(transferAmount float64, balanceAccount float64) (bool, error) {
	if transferAmount <= balanceAccount {
		return true, nil
	}
	return false, nil
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

func CreateTransaction(db *gorm.DB, amount float64, transactionType string, account Account) (Transaction, error) {

	newTransaction := Transaction{
		AccountID:       account.ID,
		Amount:          amount,
		TransactionType: transactionType,
		Description:     "transaction",
		Account:         account,
	}

	tx := db.Begin()
	if tx.Error != nil {
		log.Fatalf("Error to start transaction: %v", tx.Error)
	}

	if err := tx.Save(&newTransaction).Error; err != nil {
		tx.Rollback()
		log.Fatalf("Error to create transaction: %v", err)
		return Transaction{}, nil
	}

	if err := tx.Commit().Error; err != nil {
		log.Fatalf("Error to commit  transaction: %v", err)
	}
	return newTransaction, nil
}

func BankStatement(db *gorm.DB, accountId uint, start string, end string) (*[]Transaction, error) {
	var transactions []Transaction

	startDate, _ := helper.ToDate(start)
	endDate, _ := helper.ToDate(end)

	if err := db.Where("account_id = ? and created_at BETWEEN ? AND ?", accountId, startDate, endDate).Find(&transactions).Error; err != nil {
		log.Fatalf("Error to search transaction: %v", err)
	}
	fmt.Printf("Transactions found: %+v\n", transactions)

	return &transactions, nil
}

func Transfer(db *gorm.DB, amount float64, accountFrom Account, accountTo Account) {

	if err := db.Where("branch = ? AND number = ?", accountFrom.Branch, accountFrom.Number).Find(&accountFrom).Error; err != nil {
		log.Fatalf("Error to search transaction: %v", err)
	}

	isAble, err := IsAbleToTransfer(amount, accountFrom.Balance)
	if !isAble {
		fmt.Printf("Unable to process transfer: insufficient funds : %v\n", amount)
		return
	}

	tx := db.Begin()
	if tx.Error != nil {
		log.Fatalf("Error to start transaction: %v", tx.Error)
	}

	accountFrom.Balance -= amount
	if err := tx.Save(&accountFrom).Error; err != nil {
		tx.Rollback()
		log.Fatalf("Error to create transaction: %v", err)
	}

	if err := db.Where("branch = ? AND number = ?", accountTo.Branch, accountTo.Number).Find(&accountTo).Error; err != nil {
		log.Fatalf("Error to search transaction: %v", err)
	}

	accountTo.Balance += amount

	if err := tx.Save(&accountTo).Error; err != nil {
		tx.Rollback()
		log.Fatalf("Error to create transaction: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		log.Fatalf("Error to commit  transaction: %v", err)
	}

	_, err = CreateTransaction(db, amount, "transfer", accountTo)
	if err != nil {
		return
	}

	_, err = CreateTransaction(db, amount, "received", accountFrom)
	if err != nil {
		return
	}
}

func IsAccountActive(db *gorm.DB, account Account) (bool, error) {

	if err := db.Where("branch = ? AND number = ?", account.Branch, account.Number).Find(&account).Error; err != nil {
		log.Fatalf("Error to search account: %v", err)
	}
	if account.Status != Active {
		return false, errors.New(account.Number)
	}
	return true, nil
}
