package binding

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// StatusError 可由业务层实现，用于返回非 500 的 HTTP 状态码。
type StatusError interface {
	error
	HTTPStatus() int
}

func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	var validationErr *ValidationError
	if errors.As(err, &validationErr) {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	status := http.StatusInternalServerError
	var statusErr StatusError
	if errors.As(err, &statusErr) {
		status = statusErr.HTTPStatus()
	}

	if r != nil {
		log.Printf("%s %s -> %d: %v", r.Method, r.URL.Path, status, err)
	} else {
		log.Printf("handler error -> %d: %v", status, err)
	}

	WriteJSON(w, status, map[string]string{"error": err.Error()})
}
