package response

import (
	"encoding/json"
	"net/http"
)

// Error is an object that will be sent in http errors.
type Error struct {
	Message string `json:"message"`
}

// JSONError will write a given error message in a json object to a response writer along with the status code.
func JSONError(w http.ResponseWriter, msg string, statusCode int) {
	resp := Error{
		Message: msg,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}
