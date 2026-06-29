package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string){
	type errorMsg struct{
		Error string `json:"error"`
	}
	respondWithJson(w, code, errorMsg{Error: msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload any){
	data, err := json.Marshal(payload)
	if err != nil{
		log.Printf("Failed to Marshal JSON: %v", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
