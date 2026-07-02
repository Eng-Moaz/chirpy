package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/Eng-Moaz/chirpy/internal/auth"
	"github.com/Eng-Moaz/chirpy/internal/database"
)

func (cfg *apiConfig) HandlerRefresh(w http.ResponseWriter, r *http.Request){
	token, err := auth.GetBearerToken(r.Header)	
	if err != nil{
		respondWithError(w, 401, fmt.Sprintf("Something went wrong: %v", err))
		return
	}
	
	tokenFromDb, err := cfg.db.CheckToken(r.Context(), token)
	if err != nil || tokenFromDb.RevokedAt.Valid || tokenFromDb.ExpiresAt.Time.Before(time.Now()){
		respondWithError(w, 401, "Something went wrong, Token unauthorized")
		return
	}
	type parameters struct{
		Token string `json:"token"`
	}
	acessToken, err := auth.MakeJWT(tokenFromDb.UserID, cfg.jwt, time.Hour)
	if err != nil{
		respondWithError(w, 401, fmt.Sprintf("Something went wrong: %v", err))
		return
	}
	respondWithJson(w, 200, parameters{Token: acessToken})
}

func (cfg *apiConfig) HandlerRevoke(w http.ResponseWriter, r *http.Request){
	token, err := auth.GetBearerToken(r.Header)	
	if err != nil{
		respondWithError(w, 401, fmt.Sprintf("Something went wrong: %v", err))
		return
	}
	tokenParams := database.UpdateRefreshTokenParams{
		Token: token,
		UpdatedAt: time.Now(),
		RevokedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}
	err = cfg.db.UpdateRefreshToken(r.Context(), tokenParams)
	if err != nil{
		respondWithError(w, 401, fmt.Sprintf("Something went wrong: %v", err))
		return
	}
	respondWithJson(w, 204, struct{}{})
}
