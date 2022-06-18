package response

import (
	"encoding/json"
	"net/http"
)

// JSON will write a given object to a response writer in a json format.
func JSON(w http.ResponseWriter, obj any, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)
	if statusCode != http.StatusNoContent && obj != nil {
		json.NewEncoder(w).Encode(obj)
	}
}
