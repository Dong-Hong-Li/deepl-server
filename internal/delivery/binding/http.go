package binding

import (
	"encoding/json"
	"errors"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, err error) {
	var validationErr *ValidationError
	if errors.As(err, &validationErr) {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
}
