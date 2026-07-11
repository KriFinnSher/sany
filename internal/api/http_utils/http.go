package http_utils

import (
	"encoding/json"
	"net/http"
)

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, Error{Error: message})
}

func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func ContentType(value string) string {
	if value == "" {
		return "application/octet-stream"
	}
	return value
}

type Error struct {
	Error string `json:"error"`
}
