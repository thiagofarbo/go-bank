package main

import (
	"fmt"
	"go-bank/account"
	"go-bank/db"
	"go-bank/user"
	"go-bank/util"
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

	if action == "create" {
		userId := util.GenerateId()
		newUser := user.User{Name: strconv.Itoa(int(userId)) + "username", Age: 36, Email: "user-" + strconv.Itoa(int(userId)) + "@gmail.com"}
		newUser, _ = user.CreateUser(db.GetDB(), newUser)
		fmt.Println("User " + newUser.Name + " has been created successfully!")

		accountNumber := util.GenerateAccountNumber()

		fmt.Println("Creating account for user: " + newUser.Name)
		branch := util.GenerateAccountBranch()
		newAccount := account.Account{Branch: branch, Number: accountNumber, Balance: 0, Status: Active}

		if _, err := account.CreateAccount(db.GetDB(), newAccount, newUser); err != nil {
			fmt.Printf("Error to create an account: %v\n", err)
			return
		}
		fmt.Println("Account has been created successfully for user: " + newUser.Name)

	} else if action == "deposit" {

		deposit, _ := account.Deposit(db.GetDB(), "6662", "260456", 10.00)
		fmt.Println("Deposit made successfully for account: " + deposit.Number)

	} else if action == "withdraw" {
		withdraw, _ := account.Withdraw(db.GetDB(), "6662", "260456", 5.00)
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
	}
}
