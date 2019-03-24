package utils

import (
	"encoding/json"
	"net/http"
	"shortener/models"
)

// RespondWithError send an error to the client
func RespondWithError(w *http.ResponseWriter, status int, err models.Error) {
	(*w).Header().Add("Content-Type", "application/json")
	(*w).WriteHeader(status)
	json.NewEncoder(*w).Encode(err)
}

// RespondWithJSON send and json response to the client
func RespondWithJSON(w *http.ResponseWriter, status int, object interface{}) {
	(*w).Header().Add("Content-Type", "application/json")
	(*w).WriteHeader(status)
	json.NewEncoder(*w).Encode(object)
}

// RedirectToAnotherResource redirects client to another resource
func RedirectToAnotherResource(w *http.ResponseWriter, resource string) {
	(*w).Header().Add("Location", resource)
	(*w).WriteHeader(http.StatusMovedPermanently)
}
