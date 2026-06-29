package main

import "net/http"

func (cfg *apiConfig) handlerMetricsReset(w http.ResponseWriter, r *http.Request){
	if cfg.Platform == "dev"{
		err := cfg.db.Reset(r.Context())
		if err != nil{
			respondWithError(w, 400, "Failed to reset")
		}
		respondWithJson(w, 200, nil)
	}else{
		respondWithError(w, 400, "Forbidden")
	}
}


