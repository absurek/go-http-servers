package main

import (
	"log"
	"net/http"
	"os"

	"github.com/absurek/go-http-servers/internal/application"
	"github.com/joho/godotenv"
)

const addr = ":8080"

func main() {
	godotenv.Load()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	logger.Printf("Initializing application...")
	app, err := application.NewApplication(addr, logger)
	if err != nil {
		logger.Fatalf("ERROR: Application init: %v", err)
	}
	defer app.Close()

	go func() {
		logger.Printf("Server starting on %s", addr)

		err := app.ListenAndServe()
		if err == http.ErrServerClosed {
			logger.Println("HTTP server closed")
		} else {
			logger.Fatal("ERROR: HTTP Server:", err)
		}
	}()

	app.ListenForInterrupt()
}
