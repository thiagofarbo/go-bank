package routes

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go-bank/src/controller"
	"net/http"
)

func InitRoutes() {

	router := mux.NewRouter()

	router.HandleFunc("/users/login", controller.LoginUser).Methods("POST")
	router.HandleFunc("/users", controller.CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}", controller.GetUserById).Methods("GET")

	router.HandleFunc("/clients", controller.CreateClient).Methods("POST")
	router.HandleFunc("/clients", controller.ListClient).Methods("GET")
	router.HandleFunc("/clients/{id}", controller.GetClientById).Methods("GET")

	router.HandleFunc("/accounts", controller.CreateAccount).Methods("POST")
	router.HandleFunc("/accounts", controller.ListAccount).Methods("GET")

	router.HandleFunc("/accounts/deposit", controller.Deposit).Methods("POST")
	router.HandleFunc("/accounts/withdraw", controller.Withdraw).Methods("POST")
	router.HandleFunc("/accounts/transfer", controller.Transfer).Methods("POST")
	router.HandleFunc("/accounts/{id}/statement", controller.GetStatement).Methods("GET")
	router.HandleFunc("/accounts/{id}/loan", controller.GenerateLoan).Methods("POST")
	router.HandleFunc("/accounts/update/status", controller.UpdateStatusAccount).Methods("PUT")
	router.HandleFunc("/clients/gross-incomes", controller.AddGrossIncome).Methods("POST")
	router.HandleFunc("/hello", controller.GetHello).Methods("GET")

	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})

	http.ListenAndServe(":8080", handlers.CORS(allowedHeaders, allowedMethods, allowedOrigins)(router))
}
