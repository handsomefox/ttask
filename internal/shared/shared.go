// shared is a package with functionality that is shared in middleware and handler packages to avoid dependency cycles.
package shared

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// in case of failure return JSON { "error":"Incorrect input"} with error status code 400 Bad Request
func WriteIncorrectInputError(w http.ResponseWriter) {
	if err := WriteJSON(w, map[string]string{"error": "Incorrect input"}, http.StatusBadRequest); err != nil {
		slog.Error("Error when writing incorrect input to stream", "err", err)
	}
}

func WriteJSON(w http.ResponseWriter, v any, status int) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
