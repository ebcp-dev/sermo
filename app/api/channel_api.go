package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	utils "github.com/ebcp-dev/sermo/app/utils"
	model "github.com/ebcp-dev/sermo/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Initialize Channel API.
func (api *Api) ChannelInitialize() {
	api.initializeChannelRoutes()
}

// Defines routes.
func (api *Api) initializeChannelRoutes() {
	api.Router.HandleFunc("/api/channel", api.channelHome).Methods("GET")
	api.Router.HandleFunc("/api/channel/{id}", api.getChannel).Methods("GET")
	// Authorized routes.
	api.Router.Handle("/api/channel", api.isAuthorized(api.createChannel)).Methods("POST")
	api.Router.Handle("/api/channels", api.isAuthorized(api.getChannels)).Methods("GET")
	api.Router.Handle("/api/channel/{id}", api.isAuthorized(api.updateChannel)).Methods("PUT")
	api.Router.Handle("/api/channel/{id}", api.isAuthorized(api.deleteChannel)).Methods("DELETE")
}

// Route handlers

// Serve homepage
func (api *Api) channelHome(w http.ResponseWriter, r *http.Request) {
	current_env := os.Getenv("ENV")
	if current_env == "" {
		current_env = "dev"
	}
	fmt.Fprintln(w, "Welcome to Sermo's - Channels API")
	fmt.Fprintf(w, "ENV: %s", current_env)
}

// Retrieves channel from db using id from URL.
func (api *Api) getChannel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	ch := model.Channel{ChannelID: id}
	if err := ch.GetChannel(d.Database); err != nil {
		utils.DBNoRowsError(w, err, ch)
		return
	}
	// If channel found respond with channel object.
	utils.RespondWithJSON(w, http.StatusOK, ch)
}

// Gets list of channel with count and start variables from URL.
func (api *Api) getChannels(w http.ResponseWriter, r *http.Request) {
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
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, channel)
}

// Inserts new channel into db.
func (api *Api) createChannel(w http.ResponseWriter, r *http.Request) {
	var ch model.Channel
	// Gets JSON object from request body.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&ch); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	if err := ch.CreateChannel(d.Database); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with newly created channel.
	utils.RespondWithJSON(w, http.StatusCreated, ch)
}

// Updates channel in db using id from URL.
func (api *Api) updateChannel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	var ch model.Channel
	// Gets JSON object from request body.
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&ch); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()
	ch.ChannelID = id

	if err := ch.UpdateChannel(d.Database); err != nil {
		utils.DBNoRowsError(w, err, ch)
		return
	}
	// Respond with updated channel.
	utils.RespondWithJSON(w, http.StatusOK, ch)
}

// Deletes channel in db using id from URL.
func (api *Api) deleteChannel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Convert id string variable to int.
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}

	ch := model.Channel{ChannelID: id}
	if err := ch.DeleteChannel(d.Database); err != nil {
		utils.DBNoRowsError(w, err, ch)
		return
	}
	// Respond with success message if operation is completed.
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "channel deleted"})
}
