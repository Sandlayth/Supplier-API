package route

import (
	"encoding/json"
	"net/http"
)

// respondJSON is a helper function to respond with JSON data.
func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}