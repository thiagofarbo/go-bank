package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	account2 "go-bank/account"
	client2 "go-bank/client"
	"go-bank/db"
	"go-bank/email"
	"go-bank/grossIncome"
	"go-bank/helper"
	"go-bank/loan"
	"log"
	"net/http"
	"strconv"
)

type RequestBody struct {
	Accounts []account2.Account `json:"accounts"`
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {

	db.Connect()

	var account account2.Account

	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var branch string
	if account.Branch != "" {
		branch = account.Branch
	} else {
		branch = helper.GenerateAccountBranch()
	}

	var accountNumber string
	if account.Number != "" {
		accountNumber = account.Number
	} else {
		accountNumber = helper.GenerateAccountNumber()
	}

	var client client2.Client
	db.GetDB().First(&client, account.ClientID)
	if client.ID == 0 {
		json.NewEncoder(w).Encode("client not found!")
		return
	}
	newAccount := account2.Account{Branch: branch, Number: accountNumber, Balance: 0, Status: account2.Active}
	newAccount, _ = account2.CreateAccount(db.GetDB(), newAccount, client.ID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&newAccount)

	defer db.GetDB().Close()
	return
}

func ListAccount(w http.ResponseWriter, r *http.Request) {

	db.Connect()

	var accounts *[]account2.Account
	accounts, _ = account2.ListAccount(db.GetDB())

	w.Header().Set("Content-Type", "application/json")
	if accounts == nil || len(*accounts) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err := json.NewEncoder(w).Encode(accounts); err != nil {
		http.Error(w, "Failed to encode transactions to JSON", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	defer db.GetDB().Close()
}

func AddGrossIncome(w http.ResponseWriter, r *http.Request) {

	db.Connect()

	var accountGrossIncome grossIncome.GrossIncome
	var account account2.Account

	if err := json.NewDecoder(r.Body).Decode(&accountGrossIncome); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db.GetDB().First(&account, accountGrossIncome.AccountID)
	if account.ID == 0 {
		json.NewEncoder(w).Encode("account not found!")
		return
	}

	income := grossIncome.GrossIncome{AccountID: account.ID, Account: account, Amount: accountGrossIncome.Amount}
	income, err := grossIncome.Create(db.GetDB(), income)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&income)

	defer db.GetDB().Close()

}

func Deposit(w http.ResponseWriter, r *http.Request) {

	db.Connect()

	var account account2.Account

	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	isActive, err := account2.IsAccountActive(db.GetDB(), account)
	if !isActive {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error: account is not active for deposit: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	deposit, _ := account2.Deposit(db.GetDB(), account.Branch, account.Number, account.Balance)
	fmt.Println("Deposit made successfully for account: " + deposit.Number)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&deposit)

	defer db.GetDB().Close()
}

func Withdraw(w http.ResponseWriter, r *http.Request) {

	db.Connect()

	var account account2.Account

	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isActive, err := account2.IsAccountActive(db.GetDB(), account)
	if !isActive {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error: account is not active for withdrawals: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	withdraw, _ := account2.Withdraw(db.GetDB(), account.Branch, account.Number, account.Balance)
	fmt.Println("Withdraw made successfully for account: " + withdraw.Number)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&withdraw)

	defer db.GetDB().Close()
}

func Transfer(w http.ResponseWriter, r *http.Request) {

	db.Connect()

	var reqBody RequestBody

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	for _, value := range reqBody.Accounts {
		isActive, err := account2.IsAccountActive(db.GetDB(), value)
		if !isActive {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Fatalf("Error: account is not active for transfers: %v", err)
			return
		}
	}
	var accountFrom account2.Account
	var accountTo account2.Account
	if len(reqBody.Accounts) > 0 {

		accountFrom = reqBody.Accounts[0]
		accountTo = reqBody.Accounts[1]
		fmt.Printf("Starting transfer from account: Number=%s to account Number=%s\n", accountFrom.Number, accountTo.Number)
		account2.Transfer(db.GetDB(), accountFrom.Balance, accountFrom, accountTo)
	} else {
		http.Error(w, "No accounts provided", http.StatusBadRequest)
	}
	email.SendEmail(accountFrom, accountTo)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	defer db.GetDB().Close()
}

func GetStatement(w http.ResponseWriter, r *http.Request) {

	db.Connect()

	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	accountId := mux.Vars(r)["id"]

	u, err := strconv.ParseUint(accountId, 10, 32)
	if err != nil {
		fmt.Println("Erro ao converter string para uint:", err)
		return
	}

	transactions, err := account2.BankStatement(db.GetDB(), uint(u), start, end)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if transactions == nil || len(*transactions) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		http.Error(w, "Failed to encode transactions to JSON", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	defer db.GetDB().Close()
}

func GenerateLoan(w http.ResponseWriter, r *http.Request) {

	db.Connect()

	amount := r.URL.Query().Get("amount")
	interestRate := r.URL.Query().Get("interestRate")
	term := r.URL.Query().Get("term")
	description := r.URL.Query().Get("description")
	accountId := mux.Vars(r)["id"]

	u, err := strconv.ParseUint(accountId, 10, 32)
	if err != nil {
		fmt.Println("Error to convert string to uint:", err)
		return
	}

	loan, err := loan.GenerateLoan(db.GetDB(), uint(u), amount, interestRate, term, description+term+"%")
	if err != nil {
		http.Error(w, "Failed to generate loan", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&loan); err != nil {
		http.Error(w, "Failed to encode loan to JSON", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	defer db.GetDB().Close()
}

func UpdateStatusAccount(w http.ResponseWriter, r *http.Request) {

	db.Connect()

	var account account2.Account

	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isActive, err := account2.IsAccountActive(db.GetDB(), account)
	if !isActive {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error: account is not active for updates: %v", err)
		return
	}

	statusAccount, _ := account2.UpdateStatusAccount(db.GetDB(), account.Branch, account.Number, account2.Closed)
	if statusAccount.ID != 0 {
		fmt.Printf("Account updated successfully : %v\n", statusAccount)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&statusAccount)

	defer db.GetDB().Close()
}

func GetHello(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Thiago")
	return
}
