package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ebcp-dev/gorest-api/model"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Schema for user table
const userSchema = `
	create table if not exists users (
		id serial,
		email varchar(225) not null unique,
		password varchar(225) not null,
		createdat timestamp not null,
		updatedat timestamp not null,
		primary key (id)
	);
`

// Receives database credentials and connects to database.
func (a *App) Initialize(user, password, dbname string) {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	// Log errors.
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
	a.DB.Exec(userSchema)
}

// Starts the application.
func (a *App) Run(addr string) {
	log.Printf("Server listening on port: %s", addr)
	log.Fatal(http.ListenAndServe(":8010", a.Router))
}

// Defines routes.
func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/users", a.getUsers).Methods("GET")
	a.Router.HandleFunc("/user", a.createUser).Methods("POST")
	a.Router.HandleFunc("/user/login", a.loginUser).Methods("POST")
	a.Router.HandleFunc("/user/{id:[0-9]+}", a.getUser).Methods("GET")
	a.Router.HandleFunc("/user/{id:[0-9]+}", a.updateUser).Methods("PUT")
	a.Router.HandleFunc("/user/{id:[0-9]+}", a.deleteUser).Methods("DELETE")
}

// Route handlers

// Retrieves user from db using id from URL.
func (a *App) loginUser(w http.ResponseWriter, r *http.Request) {
	var u model.User
	// Gets JSON object from request body.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()
	// Find user in db with email and password from request body.
	log.Print(u.Password)
	if err := u.GetUserByEmail(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			// Respond with 404 if user not found in db.
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			// Respond if internal server error.
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	// Respond with user in db.
	respondWithJSON(w, http.StatusOK, u)
}

// Retrieves user from db using id from URL.
func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := strconv.Atoi(vars["id"])
	// Respond with error if id is wrong format/type.
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	u := model.User{ID: id}
	if err := u.GetUser(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			// Respond with 404 if user not found in db.
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			// Respond if internal server error.
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	// If user found respond with user object.
	respondWithJSON(w, http.StatusOK, u)
}

// Gets list of user with count and start variables from URL.
func (a *App) getUsers(w http.ResponseWriter, r *http.Request) {
	// Convert count and start string variables to int.
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	// Default and limit of count is 10.
	if count > 10 || count < 1 {
		count = 10
	}
	// Min start is 0;
	if start < 0 {
		start = 0
	}

	users, err := model.GetUsers(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, users)
}

// Inserts new user into db.
func (a *App) createUser(w http.ResponseWriter, r *http.Request) {
	var u model.User
	// Gets JSON object from request body.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	if err := u.CreateUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with newly created user.
	respondWithJSON(w, http.StatusCreated, u)
}

// Updates user in db using id from URL.
func (a *App) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := strconv.Atoi(vars["id"])
	// Respond with error if id is wrong format/type.
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var u model.User
	// Gets JSON object from request body.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()
	u.ID = id

	if err := u.UpdateUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with updated user.
	respondWithJSON(w, http.StatusOK, u)
}

// Deletes user in db using id from URL.
func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := strconv.Atoi(vars["id"])
	// Respond with error if id is wrong format/type.
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	u := model.User{ID: id}
	if err := u.DeleteUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with success message if operation is completed.
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// Helper functions

// Error message response.
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// JSON http response.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
