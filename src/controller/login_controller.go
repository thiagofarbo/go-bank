package controller

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"go-bank/src/db"
	"go-bank/src/model"
	"go-bank/src/repository"
	"net/http"
	"time"
)

var jwtKey = []byte("secret")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token,omitempty"` // token omitted if not exists
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
		//w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "Error to decode the request", http.StatusBadRequest)
		return
	}
	var user *model.User
	user, err = repository.Login(db.GetDB(), loginRequest.Username, loginRequest.Password)
	if user.ID == 0 || err != nil {
		http.Error(w, "User not found:", http.StatusNotFound)
		return
	}
	expirationTime := time.Now().Add(time.Minute * 1)

	claims := Claims{
		Username: loginRequest.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		http.Error(w, "Error to encode the token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
		Secure:   true,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Message: "User logged in successfully",
		Token:   tokenString,
	})
	db.GetDB().Close()
	return
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	//TODO - to be implemented
}
