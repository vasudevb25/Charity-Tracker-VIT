package handlers

import (
	"encoding/json"
	"my-go-backend/services"
	"net/http"
	"my-go-backend/models"
)

// TrustScoreHandler handles the POST request to calculate trust score
func TrustScoreHandler(w http.ResponseWriter, r *http.Request) {
	var ts models.TrustScore

	// Decode incoming JSON request body into the TrustScore model
	if err := json.NewDecoder(r.Body).Decode(&ts); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Calculate trust score using the service function
	score := services.CalculateTrustScore(ts)

	// Send the calculated trust score as JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"trust_score": score})
}
