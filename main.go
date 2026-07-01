package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Eng-Moaz/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	Platform string
	jwt string
}


func main(){
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil{
		log.Fatal("Failed to open a database connection")
	}
	dbQueries := database.New(db)

	jwt := os.Getenv("JWT")
	platform := os.Getenv("PLATFORM")
	apiCfg := apiConfig{db: dbQueries, Platform: platform, jwt: jwt}
	mu := http.NewServeMux()
	server := &http.Server{
		Handler: mu,
		Addr: ":8080",
	}

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	mu.Handle("/app/", fsHandler)
	mu.HandleFunc("GET /api/healthz", HandlerHealthz)
	mu.HandleFunc("GET /admin/metrics", apiCfg.HandlerMetricsWriter)
	mu.HandleFunc("POST /admin/reset", apiCfg.handlerMetricsReset)
	mu.HandleFunc("POST /api/chirps", apiCfg.HandlerChirps)
	mu.HandleFunc("POST /api/users", apiCfg.HandlerCreateUser)
	mu.HandleFunc("GET /api/chirps", apiCfg.HandlerAllChirps)
	mu.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.HandlerOneChirp)
	mu.HandleFunc("POST /api/login", apiCfg.HandlerLogin)

	err = server.ListenAndServe()
	if err != nil{
		log.Fatalf("Failed to start server: %v", err)
	}
}
