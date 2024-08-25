package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go-bank/src/db"
	"go-bank/src/helper"
	"go-bank/src/model"
	"go-bank/src/repository"
	"net/http"
)

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

	slice := []int{1, 2, 3}
	fmt.Println(slice)

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
