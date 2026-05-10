package main

import (
	"log"
)

func main() {
	app := NewApp()

	if err := app.Run(); err != nil {
		log.Fatalf("Server stopped with error: %v", err)
	}
}
