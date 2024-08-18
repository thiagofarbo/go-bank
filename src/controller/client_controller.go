package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go-bank/src/db"
	"go-bank/src/model"
	"go-bank/src/repository"
	"net/http"
)

func CreateClient(w http.ResponseWriter, r *http.Request) {
	db.Connect()

	var client model.Client

	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newClient := model.Client{Name: client.Name, Age: client.Age, Email: client.Email}
	newClient, _ = repository.CreateClient(db.GetDB(), newClient)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&newClient)

	defer db.GetDB().Close()
	return

}

func ListClient(w http.ResponseWriter, r *http.Request) {

	db.Connect()

	var clients *[]model.Client
	clients, _ = repository.ListClient(db.GetDB())

	w.Header().Set("Content-Type", "application/json")
	if clients == nil || len(*clients) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err := json.NewEncoder(w).Encode(clients); err != nil {
		http.Error(w, "Failed to encode transactions to JSON", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	defer db.GetDB().Close()
}

func GetClientById(w http.ResponseWriter, r *http.Request) {

	db.Connect()

	var client model.Client
	id := mux.Vars(r)["id"]
	db.GetDB().First(&client, id)
	if client.ID == 0 {
		json.NewEncoder(w).Encode("client not found!")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(&client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
