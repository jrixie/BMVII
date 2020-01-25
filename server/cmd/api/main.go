package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	serverAddress := "0.0.0.0:4000"
	log.Println("Server starting...")

	server := &http.Server{
		Handler:      nil,
		Addr:         serverAddress,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
