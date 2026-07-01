package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error){
	timeNow := time.Now()
	claims := &jwt.RegisteredClaims{
		Issuer: "chirpy-access",
		IssuedAt: jwt.NewNumericDate(timeNow),
		ExpiresAt: jwt.NewNumericDate(timeNow.Add(expiresIn)),
		Subject: userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(tokenSecret))
	if err != nil{
		return "", fmt.Errorf("Failed to create token")
	}
	return ss, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error){
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil{
		return uuid.Nil, fmt.Errorf("Failed to Validate")
	}else if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
		parsedID, err := uuid.Parse(claims.Subject)
		if err != nil{
			return uuid.Nil, fmt.Errorf("Failed to parse UUID")
		}
		return parsedID, nil
	}else{
		return uuid.Nil, fmt.Errorf("Failed to ValidateJWT")
	}
	
}

func GetBearerToken(headers http.Header) (string, error){
	tokenHeader := headers.Get("Authorization")
	if tokenHeader == ""{
		return "", fmt.Errorf("Failed to get Authorization header")
	}
	token := strings.TrimPrefix(tokenHeader, "Bearer ")
	return token, nil
}
