package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	app "github.com/ebcp-dev/sermo/app/utils"
	"github.com/ebcp-dev/sermo/model"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Initialize DB and routes.
func (a *App) ChannelInitialize() {
	a.initializeChannelRoutes()
}

// Defines routes.
func (a *App) initializeChannelRoutes() {
	a.Router.HandleFunc("/channel", a.channelHome).Methods("GET")
	a.Router.HandleFunc("/channel/{id}", a.getChannel).Methods("GET")
	// Authorized routes.
	a.Router.Handle("/channel", a.isAuthorized(a.createChannel)).Methods("POST")
	a.Router.Handle("/channels", a.isAuthorized(a.getChannels)).Methods("GET")
	a.Router.Handle("/channel/{id}", a.isAuthorized(a.updateChannel)).Methods("PUT")
	a.Router.Handle("/channel/{id}", a.isAuthorized(a.deleteChannel)).Methods("DELETE")
}

// Route handlers

// Serve homepage
func (a *App) channelHome(w http.ResponseWriter, r *http.Request) {
	current_env := os.Getenv("ENV")
	if current_env == "" {
		current_env = "dev"
	}
	fmt.Fprintln(w, "Welcome to Sermo's - Channels API")
	fmt.Fprintf(w, "ENV: %s", current_env)
}

// Retrieves channel from db using id from URL.
func (a *App) getChannel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	ch := model.Channel{ChannelID: id}
	if err := ch.GetChannel(d.Database); err != nil {
		switch err {
		case sql.ErrNoRows:
			// Respond with 404 if channel not found in db.
			app.RespondWithError(w, http.StatusNotFound, "Channel not found")
		default:
			// Respond if internal server error.
			app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	// If channel found respond with channel object.
	app.RespondWithJSON(w, http.StatusOK, ch)
}

// Gets list of channel with count and start variables from URL.
func (a *App) getChannels(w http.ResponseWriter, r *http.Request) {
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

	channel, err := model.GetChannels(d.Database, start, count)
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	app.RespondWithJSON(w, http.StatusOK, channel)
}

// Inserts new channel into db.
func (a *App) createChannel(w http.ResponseWriter, r *http.Request) {
	var ch model.Channel
	// Gets JSON object from request body.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&ch); err != nil {
		app.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	if err := ch.CreateChannel(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with newly created channel.
	app.RespondWithJSON(w, http.StatusCreated, ch)
}

// Updates channel in db using id from URL.
func (a *App) updateChannel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	var ch model.Channel
	// Gets JSON object from request body.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&ch); err != nil {
		app.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()
	ch.ChannelID = id

	if err := ch.UpdateChannel(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with updated channel.
	app.RespondWithJSON(w, http.StatusOK, ch)
}

// Deletes channel in db using id from URL.
func (a *App) deleteChannel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	ch := model.Channel{ChannelID: id}
	if err := ch.DeleteChannel(d.Database); err != nil {
		app.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with success message if operation is completed.
	app.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// Helper functions
