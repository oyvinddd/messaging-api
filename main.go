package main

import (
	"github.com/oyvinddd/messaging-api/internal/app"
)

func main() {
	log.Fatal(app.New().Run())
}
