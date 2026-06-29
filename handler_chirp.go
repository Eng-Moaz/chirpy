package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Eng-Moaz/chirpy/internal/database"
	"github.com/google/uuid"
)


type Chirp struct{
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}

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
		return
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
		return
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


func (cfg *apiConfig) HandlerAllChirps(w http.ResponseWriter, r *http.Request){
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil{
		respondWithError(w, 400, "Failed to get chirps")
		return
	}
	chirpsResponse := make([]Chirp, 0)

	for _, chirp := range chirps{
		curChirp := Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserId: chirp.UserID,
		}
		chirpsResponse = append(chirpsResponse, curChirp)
	}
	respondWithJson(w, 200, chirpsResponse)
}


func (cfg *apiConfig) HandlerOneChirp(w http.ResponseWriter, r *http.Request){
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil{
		respondWithError(w, 400, "Something went wrong")
	}
	chirp, err := cfg.db.GetChirpById(r.Context(),chirpID) 

	if err != nil{
		respondWithError(w, 404, "Chirp not found")
	}

	chirpsResponse := Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserId: chirp.UserID,
		}
	respondWithJson(w, 200, chirpsResponse)
}
