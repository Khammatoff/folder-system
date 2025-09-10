package handler

import (
	"encoding/json"
	"net/http"
)

// WriteJSONError пишет JSON-ответ с ошибкой и соответствующим статусом HTTP.
func WriteJSONError(w http.ResponseWriter, status int, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": errMsg})
}
