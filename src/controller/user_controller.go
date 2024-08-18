package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go-bank/src/db"
	"go-bank/src/helper"
	"go-bank/src/model"
	"go-bank/src/repository"
	"net/http"
)

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token,omitempty"` // token omitido se não existir
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginUser(w http.ResponseWriter, r *http.Request) {

	db.Connect()

	var loginRequest LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, "Erro ao decodificar a requisição", http.StatusBadRequest)
		return
	}
	var user *model.User
	user, err = repository.Login(db.GetDB(), loginRequest.Username, loginRequest.Password)
	if user.ID == 0 || err != nil {
		http.Error(w, "User not found:", http.StatusBadRequest)
		return
	}

	token := "token"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Message: "Login realizado com sucesso",
		Token:   token,
	})
	db.GetDB().Close()
	return
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	db.Connect()

	var user model.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashPassword, err := helper.EncryptPassword(user.Password)
	if err != nil {
		return
	}

	newUser := model.User{Username: user.Username, Password: string(hashPassword), Email: user.Email}
	newUser, _ = repository.CreateUser(db.GetDB(), newUser)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&user)

	defer db.GetDB().Close()
	return
}

func GetUserById(w http.ResponseWriter, r *http.Request) {

	db.Connect()

	var user model.User
	id := mux.Vars(r)["id"]
	db.GetDB().First(&user, id)
	if user.ID == 0 {
		//http.Error(w, "user not found", http.StatusNotFound)
		json.NewEncoder(w).Encode(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
