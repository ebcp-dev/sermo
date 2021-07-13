package utils

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
)

// Responds with error if no rows were found in the query.
func DBNoRowsError(w http.ResponseWriter, err error, obj interface{}) {
	// Format message for specific object type.
	objType := fmt.Sprintf("%T", obj)
	errMsg := strings.TrimLeft(objType+" not found", "model.")

	switch err {
	case sql.ErrNoRows:
		// Respond with 404 if object not found in db.
		RespondWithError(w, http.StatusNotFound, errMsg)
	default:
		// Respond if internal server error.
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
}
