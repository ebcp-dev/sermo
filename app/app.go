package app

import (
	"log"
	"net/http"

	"github.com/ebcp-dev/sermo/app/api"
)

// References Api struct in api package.
var a api.Api

type App struct{}

// Initialize DB and routes.
func (app *App) Initialize() {
	a.InitializeAPI()
}

// Starts the application.
func (app *App) Run(addr string) {
	log.Printf("Server listening on port: %s", addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
