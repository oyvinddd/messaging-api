package main

import (
	"log"
	"github.com/oyvinddd/messaging-api/internal/app"
)

func main() {
	cfg := *app.NewDefaultConfig()
	log.Fatal(app.New(cfg).Run())
}
