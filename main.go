package main

import (
	"log"
	"net/http"
)

func main(){
	mu := http.NewServeMux()
	server := http.Server{
		Handler: mu,
		Addr: ":8080",
	}
	mu.Handle("/", http.FileServer(http.Dir(".")))
	err := server.ListenAndServe()
	if err != nil{
		log.Fatal("Failed to start server")
	}
}
