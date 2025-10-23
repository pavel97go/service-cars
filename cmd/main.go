package main

import (
	"context"
	"log"

	"github.com/pavel97go/service-cars/internal/app"
)

func main() {
	if err := app.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
