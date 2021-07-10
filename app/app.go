package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	utils "github.com/ebcp-dev/sermo/app/utils"
	"github.com/ebcp-dev/sermo/db"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

// References DB struct in db.go.
var d db.DB

type App struct {
	Router *mux.Router
}

// Initialize DB and routes.
func (a *App) Initialize() {
	// Find and read the config file.
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	// Default is ENV=dev.
	db_user := viper.GetString("APP_DB_USERNAME")
	db_pass := viper.GetString("APP_DB_PASSWORD")
	db_host := viper.GetString("APP_DB_HOST")
	db_name := viper.GetString("APP_DB_NAME")
	// Production env variables.
	if os.Getenv("ENV") == "prod" {
		db_user = os.Getenv("PROD_DB_USERNAME")
		db_pass = os.Getenv("PROD_DB_PASSWORD")
		db_host = os.Getenv("PROD_DB_HOST")
		db_name = os.Getenv("PROD_DB_NAME")
	}

	// Receives database credentials and connects to database.
	d.Initialize(db_user, db_pass, db_host, db_name)

	// Initialize mux router.
	a.Router = mux.NewRouter()

	// Handle home page.
	a.Router.HandleFunc("/", homePage)

	// Initialize other app routes.
	a.UserInitialize()
	a.ChannelInitialize()
	a.SignalInitialize()
}

// Starts the application.
func (a *App) Run(addr string) {
	log.Printf("Server listening on port: %s", addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

// Serve homepage
func homePage(w http.ResponseWriter, r *http.Request) {
	current_env := os.Getenv("ENV")
	if current_env == "" {
		current_env = "dev"
	}
	// Show environment.
	fmt.Fprintln(w, "Welcome to Sermo - API")
	fmt.Fprintf(w, "ENV: %s", current_env)
}

// Authorization middleware
func (a *App) isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if request has "Token" header.
		authorizationHeader := r.Header["Token"]
		if !utils.ValidateToken(strings.Join(authorizationHeader, "")) {
			utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		} else {
			endpoint(w, r)
		}
	})
}
