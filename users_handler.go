package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Eng-Moaz/chirpy/internal/auth"
	"github.com/Eng-Moaz/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct{
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
	Token string `json:"token"`
}

type UserReceived struct{
	Email string `json:"email"`
	Password string `json:"password"`
	ExpiresInSeconds *int `json:"expires_in_seconds"`
}

func (cfg *apiConfig)HandlerCreateUser(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	params := UserReceived{}
	err := decoder.Decode(&params)
	if err != nil{
		respondWithError(w, 400, "Something went wrong")
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil{
		respondWithError(w, 400, "Failed to hashedPassword")
	}
	userParams := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email: params.Email,
		HashedPassword: hashedPassword,
	}
	
	user, err := cfg.db.CreateUser(r.Context(), userParams)
	if err != nil{
		respondWithError(w, 400, "Failed to create user")
	}
	userResp := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}
	respondWithJson(w, 201, userResp)
}

func (cfg *apiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	params := UserReceived{}
	err := decoder.Decode(&params)
	if err != nil{
		respondWithError(w, 400, "Something went wrong")
		return
	}

	userFromDB, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil{
		respondWithError(w, 400, "Failed to retrieve user")
		return
	}

	hashedPasswordCorrect := userFromDB.HashedPassword
	ok, err := auth.CheckPasswordHash(params.Password, hashedPasswordCorrect)
	if err != nil{
		respondWithError(w, 400, "Something went wrong")
		return
	}

	if !ok{
		respondWithError(w, 401, "Unauthorized")
		return
	}

	expiresInSeconds := time.Hour * 1
	if params.ExpiresInSeconds != nil && *params.ExpiresInSeconds >= (60 * 60){
		expiresInSeconds = time.Duration(*params.ExpiresInSeconds)
	}
	secretString, err := auth.MakeJWT(userFromDB.ID, cfg.jwt, expiresInSeconds)
	if err != nil{
		respondWithError(w, 400, "Something went wrong")
		return
	}

	userResp := User{
		ID: userFromDB.ID,
		CreatedAt: userFromDB.CreatedAt,
		UpdatedAt: userFromDB.UpdatedAt,
		Email: userFromDB.Email,
		Token: secretString,
	}
	respondWithJson(w, 200, userResp)
}
