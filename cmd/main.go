package main

import (
	"context"
	"log"

	"github.com/pavel97go/service-cars/internal/app"
)

func main() {
	ctx := context.Background()
	if err := app.Run(ctx); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
