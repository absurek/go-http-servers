package main

import (
	"log"
	"os"

	"github.com/absurek/go-http-servers/internal/application"
	"github.com/joho/godotenv"
)

const addr = ":8080"

func main() {
	godotenv.Load()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app, err := application.NewApplication(addr, logger)
	if err != nil {
		logger.Fatalf("ERROR: Application init: %v", err)
	}
	defer app.Close()

	app.ListenAndServe()
}
