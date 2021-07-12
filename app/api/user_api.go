package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/ebcp-dev/sermo/app/auth"
	utils "github.com/ebcp-dev/sermo/app/utils"
	model "github.com/ebcp-dev/sermo/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Initialize User API.
func (api *Api) UserInitialize() {
	api.initializeUserRoutes()
}

// Defines routes.
func (api *Api) initializeUserRoutes() {
	api.Router.HandleFunc("/user", api.userHome).Methods("GET")
	api.Router.HandleFunc("/user", api.createUser).Methods("POST")
	api.Router.HandleFunc("/user/login", api.loginUser).Methods("POST")
	// Authorized routes.
	api.Router.Handle("/user/{id}", api.isAuthorized(api.getUser)).Methods("GET")
	api.Router.Handle("/users", api.isAuthorized(api.getUsers)).Methods("GET")
	api.Router.Handle("/user/{id}", api.isAuthorized(api.updateUser)).Methods("PUT")
	api.Router.Handle("/user/{id}", api.isAuthorized(api.deleteUser)).Methods("DELETE")
}

// Route handlers

// Serve homepage
func (api *Api) userHome(w http.ResponseWriter, r *http.Request) {
	current_env := os.Getenv("ENV")
	if current_env == "" {
		current_env = "dev"
	}
	fmt.Fprintln(w, "Welcome to Sermo's - Users API")
	fmt.Fprintf(w, "ENV: %s", current_env)
}

// Retrieves user from db using id from URL.
func (api *Api) loginUser(w http.ResponseWriter, r *http.Request) {
	var u model.User
	// Gets JSON object from request body.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	passwordInput := u.Password

	defer r.Body.Close()
	// Find user in db with email from request body.
	if err := u.GetUserByEmail(d.Database); err != nil {
		switch err {
		case sql.ErrNoRows:
			// Respond with 404 if user not found in db.
			utils.RespondWithError(w, http.StatusNotFound, "User not found.")
			return
		default:
			// Respond if internal server error.
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
	}
	if !auth.ComparePasswords(u.Password, []byte(passwordInput)) {
		// Respond with 401 if hashed passwords don't match.
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid login.")
		return
	}
	// Generate and send token to client with response header.
	validToken, err := auth.GenerateJWT()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Add("Token", validToken)
	// Respond with user in db.
	utils.RespondWithJSON(w, http.StatusOK, u)
}

// Retrieves user from db using id from URL.
func (api *Api) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	u := model.User{UserID: id}
	if err := u.GetUser(d.Database); err != nil {
		switch err {
		case sql.ErrNoRows:
			// Respond with 404 if user not found in db.
			utils.RespondWithError(w, http.StatusNotFound, "User not found")
		default:
			// Respond if internal server error.
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	// If user found respond with user object.
	utils.RespondWithJSON(w, http.StatusOK, u)
}

// Gets list of user with count and start variables from URL.
func (api *Api) getUsers(w http.ResponseWriter, r *http.Request) {
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

	users, err := model.GetUsers(d.Database, start, count)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, users)
}

// Inserts new user into db.
func (api *Api) createUser(w http.ResponseWriter, r *http.Request) {
	var u model.User
	// Gets JSON object from request body.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	// Hash password.
	u.Password = auth.HashAndSalt([]byte(u.Password))
	defer r.Body.Close()

	if err := u.CreateUser(d.Database); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with newly created user.
	utils.RespondWithJSON(w, http.StatusCreated, u)
}

// Updates user in db using id from URL.
func (api *Api) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	var u model.User
	// Gets JSON object from request body.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()
	u.UserID = id
	// Hash password.
	u.Password = auth.HashAndSalt([]byte(u.Password))

	if err := u.UpdateUser(d.Database); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with updated user.
	utils.RespondWithJSON(w, http.StatusOK, u)
}

// Deletes user in db using id from URL.
func (api *Api) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	u := model.User{UserID: id}
	if err := u.DeleteUser(d.Database); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with success message if operation is completed.
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "user deleted"})
}
