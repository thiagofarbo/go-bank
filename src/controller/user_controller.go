package controller

import (
	"encoding/json"
	"go-bank/db"
	"go-bank/helper"
	"go-bank/login"
	user2 "go-bank/user"
	"net/http"
)

func LoginUser(w http.ResponseWriter, r *http.Request) {

	db.Connect()

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

	token := "token"

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
