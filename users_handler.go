package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
	RefreshToken string `json:"refresh_token"`
}

type UserReceived struct{
	Email string `json:"email"`
	Password string `json:"password"`
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

func (cfg *apiConfig) createRefreshToken(r *http.Request, userID uuid.UUID) (string, error){
	timeNow := time.Now()
	params := database.CreateTokenParams{
		Token: auth.MakeRefreshToken(),
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
		UserID: userID,
		ExpiresAt: sql.NullTime{Time: timeNow.AddDate(0,0,60), Valid: true},
		RevokedAt: sql.NullTime{Time: time.Time{}, Valid: false},
	}
	token, err := cfg.db.CreateToken(r.Context(), params)
	if err != nil{
		return "", fmt.Errorf("Failed to create token")
	}
	return token.Token, nil
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
	secretString, err := auth.MakeJWT(userFromDB.ID, cfg.jwt, expiresInSeconds)
	if err != nil{
		respondWithError(w, 400, "Something went wrong")
		return
	}
	refreshToken, err := cfg.createRefreshToken(r, userFromDB.ID)
	if err != nil{
		respondWithError(w, 400, "Failed to create token")
		return
	}

	userResp := User{
		ID: userFromDB.ID,
		CreatedAt: userFromDB.CreatedAt,
		UpdatedAt: userFromDB.UpdatedAt,
		Email: userFromDB.Email,
		Token: secretString,
		RefreshToken: refreshToken,
	}
	respondWithJson(w, 200, userResp)
}
