package main

import (
	"net/http"
	"time"

	"boilermakevii/api/internal/mongo"
	"boilermakevii/api/internal/router"

	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

func main() {
	setupLogging()
	defer mongo.Close()

	serverAddress := "0.0.0.0:4000"
	log.Println("Server starting...")

	apiRouter := router.NewRouter()
	handler := cors.Default().Handler(apiRouter)

	server := &http.Server{
		Handler:      handler,
		Addr:         serverAddress,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}

func setupLogging() {
	logFormatter := new(log.TextFormatter)
	logFormatter.TimestampFormat = "2006-01-02 15:04:05"
	logFormatter.FullTimestamp = true
	log.SetFormatter(logFormatter)
}
