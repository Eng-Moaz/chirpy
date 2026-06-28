package main

import (
	"encoding/json"
	"net/http"
)

func HandlerValidateChirp(w http.ResponseWriter, r *http.Request){
	type parameters struct{
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil{
		respondWithError(w, 400, "Something went wrong")
	}
	var Valid bool
	if len(params.Body) > 140{
		Valid = false
	}else{
		Valid = true
	}
	if !Valid{
		respondWithError(w, 400, "Chirp is too long")
	}else{
		type okParameters struct{
			IsValid bool `json:"valid"`
		}
		respondWithJson(w, 200, okParameters{IsValid: Valid})
	}
}
