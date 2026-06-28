package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	_ "github.com/lib/pq"
	"github.com/Eng-Moaz/chirpy/internal/databse"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	database *database.Queries
}


func main(){
	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil{
		log.Fatal("Failed to open a database connection")
	}

	apiCfg := apiConfig{database: db}
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

	err = server.ListenAndServe()
	if err != nil{
		log.Fatalf("Failed to start server: %v", err)
	}
}
