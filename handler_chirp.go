package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Eng-Moaz/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) HandlerChirps(w http.ResponseWriter, r *http.Request){
	type parameters struct{
		Body string `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}
	
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil{
		respondWithError(w, 400, "Something went wrong")
	}
	if len(params.Body) > 140{
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	cleanedBody := cleanProfanity(params.Body)
	chirpParams := database.CreateChirpParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body: cleanedBody,
		UserID: params.UserId,
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), chirpParams)
	if err != nil{
		respondWithError(w, 400, "Failed to create chirp")
	}
	type Chirp struct{
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body string `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}
	respChirp := Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserId: chirp.UserID,
	}
	respondWithJson(w, 201, respChirp)
}
