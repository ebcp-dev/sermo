package main

import (
	"log"
	"os"

	"github.com/ebcp-dev/gorest-api/app"
)

func main() {
	// Log current environment.
	current_env := os.Getenv("ENV")
	if current_env == "" {
		current_env = "dev"
	}
	log.Println("ENV: " + current_env)

	a := app.App{}

	a.Initialize()
	a.Run(":8010")
}
