package main

import (
	"encoding/json"
	"net/http"
)

func CustomJsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, `{error: failed to encode response}`, http.StatusInternalServerError)
	}
}
