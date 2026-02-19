package main

import (
	"github.com/dontpanicw/SalesTracker/internal/app"
	"log"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("Failed to run app: %v", err)
	}
}
