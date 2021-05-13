package app

import (
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	app "github.com/ebcp-dev/gorest-api/app/utils"
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
	a.DataInitialize()
}

// Starts the application.
func (a *App) Run(addr string) {
	log.Printf("Server listening on port: %s", addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

// Authorization middleware
func (a *App) isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if request has "Token" header.
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				// Check if token is valid based on private `mySigningKey`.
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					app.RespondWithError(w, http.StatusInternalServerError, "There was error with signing the token.")
				}
				return mySigningKey, nil
			})

			if err != nil {
				app.RespondWithError(w, http.StatusInternalServerError, err.Error())
			}
			// Serve endpoint if token is valid.
			if token.Valid {
				endpoint(w, r)
			}
		} else {
			app.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		}
	})
}
