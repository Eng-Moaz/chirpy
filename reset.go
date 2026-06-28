package main

import "net/http"

func (cfg *apiConfig) handlerMetricsReset(w http.ResponseWriter, r *http.Request){
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")	
	w.WriteHeader(200)
	cfg.fileserverHits.Store(0)
}


