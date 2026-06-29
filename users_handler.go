package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Eng-Moaz/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig)HandlerCreateUser(w http.ResponseWriter, r *http.Request){
	type parameters struct{
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil{
		respondWithError(w, 400, "Something went wrong")
	}
	userParams := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email: params.Email,
	}
	
	user, err := cfg.db.CreateUser(r.Context(), userParams)
	if err != nil{
		respondWithError(w, 400, "Failed to create user")
	}
	type User struct{
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
	}
	userResp := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}
	respondWithJson(w, 201, userResp)
}
