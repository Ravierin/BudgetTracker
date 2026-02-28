package handler

import (
	"github.com/Ravierin/BudgetTracker/backend/internal/model"
	"github.com/Ravierin/BudgetTracker/backend/internal/service"
	"encoding/json"
	"log"
	"net/http"
)

type APIKeyHandler struct {
	service *service.APIKeyService
}

func NewAPIKeyHandler(service *service.APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{service: service}
}

func (h *APIKeyHandler) GetAPIKeys(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	apiKeys, err := h.service.GetAllAPIKeys(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Don't expose actual secrets in response
	for i := range apiKeys {
		if apiKeys[i].APIKey != "" {
			apiKeys[i].APIKey = maskString(apiKeys[i].APIKey)
		}
		if apiKeys[i].APISecret != "" {
			apiKeys[i].APISecret = maskString(apiKeys[i].APISecret)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apiKeys)
}

func (h *APIKeyHandler) SaveAPIKeys(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var keys []model.APIKey
	if err := json.NewDecoder(r.Body).Decode(&keys); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("[API] Saving %d API keys: %v", len(keys), keys)

	if len(keys) == 0 {
		http.Error(w, "No API keys provided", http.StatusBadRequest)
		return
	}

	for _, key := range keys {
		log.Printf("[API] Saving key for %s: apiKey=%s... (len=%d)", key.Exchange, key.APIKey[:8], len(key.APIKey))
		if err := h.service.SaveAPIKey(ctx, &key); err != nil {
			http.Error(w, "Failed to save key for "+key.Exchange+": "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "saved"})
}

func maskString(s string) string {
	if len(s) <= 8 {
		return "****"
	}
	return s[:4] + "****" + s[len(s)-4:]
}
