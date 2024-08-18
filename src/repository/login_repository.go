package repository

import (
	"fmt"
	"github.com/jinzhu/gorm"
	user2 "go-bank/src/model"
	"golang.org/x/crypto/bcrypt"
)

func Login(db *gorm.DB, username string, password string) (*user2.User, error) {
	// Buscando o usuário no banco de dados
	var user user2.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		fmt.Printf("User not found : %v\n", username)
		return &user2.User{}, nil
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Printf("Wrong username or password: %v\n", username)
		return &user2.User{}, nil
	}
	return &user, nil
}
