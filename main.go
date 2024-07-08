package main

import (
	"fmt"
	"go-bank/account"
	"go-bank/client"
	"go-bank/db"
	"go-bank/grossIncome"
	"go-bank/helper"
	"go-bank/user"
	"os"
	"strconv"
)

const (
	Active    = "active"
	Inactive  = "inactive"
	Frozen    = "frozen"
	Closed    = "closed"
	Blocked   = "blocked"
	Suspended = "suspended"
)

func main() {
	db.Connect()

	action := os.Args[1]
	fmt.Println("Action:", action)

	if action == "transfer" {
		account.Transfer(db.GetDB(), 10, account.Account{Branch: "5255", Number: "554601"}, account.Account{Branch: "2236", Number: "568387"})
	}

	if action == "create-user" {
		userId := helper.GenerateId()
		newUser := user.User{Username: strconv.Itoa(int(userId)) + "username", Password: "1234", Email: "user-" + strconv.Itoa(int(userId)) + "@gmail.com"}
		newUser, _ = user.CreateUser(db.GetDB(), newUser)
		fmt.Println("User " + newUser.Username + " has been created successfully!")
	}

	if action == "create-client" {
		userId := helper.GenerateId()
		newUser := client.Client{Name: strconv.Itoa(int(userId)) + "username", Age: 36, Email: "user-" + strconv.Itoa(int(userId)) + "@gmail.com"}
		newUser, _ = client.CreateUser(db.GetDB(), newUser)
		fmt.Println("User " + newUser.Name + " has been created successfully!")

		accountNumber := helper.GenerateAccountNumber()

		fmt.Println("Creating account for user: " + newUser.Name)
		branch := helper.GenerateAccountBranch()
		newAccount := account.Account{Branch: branch, Number: accountNumber, Balance: 0, Status: Active}

		createAccount, err := account.CreateAccount(db.GetDB(), newAccount, newUser)
		if err != nil {
			fmt.Printf("Error to create an account: %v\n", err)
			return
		}
		fmt.Println("Account has been created successfully for user: " + newUser.Name)

		income := grossIncome.GrossIncome{AccountID: createAccount.ID, Account: createAccount, Amount: 25000}
		grossIncome.Create(db.GetDB(), income)

	} else if action == "deposit" {

		deposit, _ := account.Deposit(db.GetDB(), "2236", "568387", 100.00)
		fmt.Println("Deposit made successfully for account: " + deposit.Number)

	} else if action == "withdraw" {
		withdraw, _ := account.Withdraw(db.GetDB(), "5255", "554601", 5.00)
		if withdraw.ID == 0 {
			fmt.Println("Error to withdraw balance" + withdraw.Number)
		} else {
			fmt.Println("Withdraw made successfully for account: " + withdraw.Number)
		}
	} else if action == "update" {
		statusAccount, _ := account.UpdateStatusAccount(db.GetDB(), "6662", "260456", Closed)
		if statusAccount.ID != 0 {
			fmt.Printf("Account updated successfully : %v\n", statusAccount)
		}
	} else if action == "statement" {
		account.BankStatement(db.GetDB(), "2024-07-01 00:00:00", "2024-07-08 23:59:59")
	}
}
