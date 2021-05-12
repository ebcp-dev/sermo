package app

import (
	"log"
	"net/http"
	"os"

	"github.com/ebcp-dev/gorest-api/db"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// References DB struct in db.go.
var d db.DB

type App struct {
	Router *mux.Router
}

// Initialize DB and routes.
func (a *App) Initialize() {
	// Receives database credentials and connects to database.
	d.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"))

	a.Router = mux.NewRouter()
	a.UserInitialize()
}

// Starts the application.
func (a *App) Run(addr string) {
	log.Printf("Server listening on port: %s", addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
