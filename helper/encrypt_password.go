package helper

import (
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) ([]byte, error) {

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}
	return hashPassword, nil
}
