package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	account2 "go-bank/account"
	client2 "go-bank/client"
	"go-bank/db"
	"go-bank/grossIncome"
	"go-bank/helper"
	"go-bank/loan"
	"go-bank/login"
	user2 "go-bank/user"
	"log"
	"net/http"
	"strconv"
)

type RequestBody struct {
	Accounts []account2.Account `json:"accounts"`
}

func LoginUser(w http.ResponseWriter, r *http.Request) {

	db.Connect()

	// Decodificando a requisição JSON
	var loginRequest login.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, "Erro ao decodificar a requisição", http.StatusBadRequest)
		return
	}
	var user *user2.User
	user, err = login.Login(db.GetDB(), loginRequest.Username, loginRequest.Password)
	if user.ID == 0 || err != nil {
		http.Error(w, "User not found:", http.StatusBadRequest)
		return
	}

	// Gerando um token de autenticação (substitua com sua lógica de geração de token)
	token := "exemploDeToken" // Aqui você pode gerar um token JWT, por exemplo

	// Respondendo com sucesso
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(login.LoginResponse{
		Message: "Login realizado com sucesso",
		Token:   token,
	})
	db.GetDB().Close()
	return
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	db.Connect()

	var user user2.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashPassword, err := helper.EncryptPassword(user.Password)
	if err != nil {
		return
	}

	newUser := user2.User{Username: user.Username, Password: string(hashPassword), Email: user.Email}
	newUser, _ = user2.Create(db.GetDB(), newUser)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&user)

	defer db.GetDB().Close()
	return
}

func CreateClient(w http.ResponseWriter, r *http.Request) {
	db.Connect()

	var client client2.Client

	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newClient := client2.Client{Name: client.Name, Age: client.Age, Email: client.Email}
	newClient, _ = client2.Create(db.GetDB(), newClient)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&newClient)

	defer db.GetDB().Close()
	return

}

func CreateAccount(w http.ResponseWriter, r *http.Request) {

	db.Connect()

	var account account2.Account

	accountNumber := helper.GenerateAccountNumber()

	branch := helper.GenerateAccountBranch()

	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
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

func GetClientById(w http.ResponseWriter, r *http.Request) (*client2.Client, error) {

	db.Connect()

	var client client2.Client
	id := mux.Vars(r)["id"]
	db.GetDB().First(&client, id)
	if client.ID == 0 {
		json.NewEncoder(w).Encode("client not found!")
		return &client2.Client{}, errors.New("client not found!")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(&client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return &client2.Client{}, err
	}
	return &client, err
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

	if len(reqBody.Accounts) > 0 {

		accountFrom := reqBody.Accounts[0]
		accountTo := reqBody.Accounts[1]
		fmt.Printf("Starting transfer from account: Number=%s to account Number=%s\n", accountFrom.Number, accountTo.Number)
		account2.Transfer(db.GetDB(), accountFrom.Balance, account2.Account{Branch: accountFrom.Branch, Number: accountFrom.Number}, account2.Account{Branch: accountTo.Branch, Number: accountTo.Number})
	} else {
		http.Error(w, "No accounts provided", http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)

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

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/users/login", LoginUser).Methods("POST")
	router.HandleFunc("/users", CreateUser).Methods("POST")
	router.HandleFunc("/clients", CreateClient).Methods("POST")

	router.HandleFunc("/clients/{id}", func(w http.ResponseWriter, r *http.Request) {
		user, err := GetClientById(w, r)
		if err == nil {
			json.NewEncoder(w).Encode(user)
		}
	}).Methods("GET")

	router.HandleFunc("/accounts", CreateAccount).Methods("POST")
	router.HandleFunc("/accounts/deposit", Deposit).Methods("POST")
	router.HandleFunc("/accounts/withdraw", Withdraw).Methods("POST")
	router.HandleFunc("/accounts/transfer", Transfer).Methods("POST")
	router.HandleFunc("/accounts/{id}/statement", GetStatement).Methods("GET")
	router.HandleFunc("/accounts/{id}/loan", GenerateLoan).Methods("POST")
	router.HandleFunc("/accounts/update/status", UpdateStatusAccount).Methods("PUT")
	router.HandleFunc("/clients/gross-incomes", AddGrossIncome).Methods("POST")
	router.HandleFunc("/hello", GetHello).Methods("GET")

	// Configuração do CORS
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})

	// Iniciar o servidor com as configurações de CORS
	http.ListenAndServe(":8080", handlers.CORS(allowedHeaders, allowedMethods, allowedOrigins)(router))
}
