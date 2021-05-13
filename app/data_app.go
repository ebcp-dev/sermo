package app

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	app "github.com/ebcp-dev/gorest-api/app/utils"
	"github.com/ebcp-dev/gorest-api/model"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Initialize DB and routes.
func (a *App) DataInitialize() {
	a.initializeDataRoutes()
}

// Defines routes.
func (a *App) initializeDataRoutes() {
	a.Router.HandleFunc("/data/{id}", a.getData).Methods("GET")
	// Authorized routes.
	a.Router.Handle("/data", a.isAuthorized(a.createData)).Methods("POST")
	a.Router.Handle("/data", a.isAuthorized(a.getMultipleData)).Methods("GET")
	a.Router.Handle("/data/{id}", a.isAuthorized(a.updateData)).Methods("PUT")
	a.Router.Handle("/data/{id}", a.isAuthorized(a.deleteData)).Methods("DELETE")
}

// Route handlers

// Retrieves data from db using id from URL.
func (a *App) getData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	dt := model.Data{ID: id}
	if err := dt.GetData(d.Database); err != nil {
		switch err {
		case sql.ErrNoRows:
			// Respond with 404 if data not found in db.
			app.RespondWithError(w, http.StatusNotFound, "Data not found")
		default:
			// Respond if internal server error.
			app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	// If data found respond with data object.
	app.RespondWithJSON(w, http.StatusOK, dt)
}

// Gets list of data with count and start variables from URL.
func (a *App) getMultipleData(w http.ResponseWriter, r *http.Request) {
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

	data, err := model.GetMultipleData(d.Database, start, count)
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	app.RespondWithJSON(w, http.StatusOK, data)
}

// Inserts new data into db.
func (a *App) createData(w http.ResponseWriter, r *http.Request) {
	var dt model.Data
	// Gets JSON object from request body.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&dt); err != nil {
		app.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	if err := dt.CreateData(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with newly created data.
	app.RespondWithJSON(w, http.StatusCreated, dt)
}

// Updates data in db using id from URL.
func (a *App) updateData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	var dt model.Data
	// Gets JSON object from request body.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&dt); err != nil {
		app.RespondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}

	defer r.Body.Close()
	dt.ID = id

	if err := dt.UpdateData(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with updated data.
	app.RespondWithJSON(w, http.StatusOK, dt)
}

// Deletes data in db using id from URL.
func (a *App) deleteData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	dt := model.Data{ID: id}
	if err := dt.DeleteData(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with success message if operation is completed.
	app.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// Helper functions
