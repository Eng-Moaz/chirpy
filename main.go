package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}


func main(){
	apiCfg := apiConfig{}
	mu := http.NewServeMux()
	server := &http.Server{
		Handler: mu,
		Addr: ":8080",
	}

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	mu.Handle("/app/", fsHandler)
	mu.HandleFunc("GET /api/healthz", HandlerHealthz)
	mu.HandleFunc("GET /admin/metrics", apiCfg.handlerMetricsWriter)
	mu.HandleFunc("POST /admin/reset", apiCfg.handlerMetricsReset)
	mu.HandleFunc("POST /api/validate_chirp", HandlerValidateChirp)

	err := server.ListenAndServe()
	if err != nil{
		log.Fatalf("Failed to start server: %v", err)
	}
}
