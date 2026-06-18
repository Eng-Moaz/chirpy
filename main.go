package main

import (
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request){
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func main(){
	mu := http.NewServeMux()
	server := &http.Server{
		Handler: mu,
		Addr: ":8080",
	}
	mu.Handle("/app/", http.StripPrefix("/app",http.FileServer(http.Dir("."))))
	mu.HandleFunc("/healthz", handler)
	err := server.ListenAndServe()
	if err != nil{
		log.Fatal("Failed to start server")
	}
}
