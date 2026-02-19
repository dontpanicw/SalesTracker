package main

import (
	"log"

	"github.com/yourusername/analytics-service/internal/app"
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
