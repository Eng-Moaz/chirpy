package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler{
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request){
			cfg.fileserverHits.Add(1)
			next.ServeHTTP(w, r)
		},
	)
}

func (cfg *apiConfig) handlerMetricsWriter(w http.ResponseWriter, r *http.Request){
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")	
	w.WriteHeader(200)
	serverHits := cfg.fileserverHits.Load()
	msgToWrite := fmt.Sprintf("Hits: %v", serverHits)
	w.Write([]byte(msgToWrite))
}

func (cfg *apiConfig) handlerMetricsReset(w http.ResponseWriter, r *http.Request){
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")	
	w.WriteHeader(200)
	cfg.fileserverHits.Store(0)
}

func handler(w http.ResponseWriter, r *http.Request){
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func main(){
	cfg := apiConfig{}
	mu := http.NewServeMux()
	server := &http.Server{
		Handler: mu,
		Addr: ":8080",
	}

	mu.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app",http.FileServer(http.Dir(".")))))
	mu.HandleFunc("GET /healthz", handler)
	mu.HandleFunc("GET /metrics", cfg.handlerMetricsWriter)
	mu.HandleFunc("POST /reset", cfg.handlerMetricsReset)

	err := server.ListenAndServe()
	if err != nil{
		log.Fatalf("Failed to start server: %v", err)
	}
}
