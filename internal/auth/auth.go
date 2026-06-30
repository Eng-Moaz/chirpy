package auth

import (
	"fmt"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error){
	hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil{
		return "", fmt.Errorf("Failed to hashedPassword")
	}
	return hashedPassword, nil
}

func CheckPasswordHash(password, hash string) (bool, error){
	ok, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil{
		return false, fmt.Errorf("Failed to CheckHash")
	}
	return ok, nil
}
