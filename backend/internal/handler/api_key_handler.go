package handler

import (
	"BudgetTracker/backend/internal/model"
	"BudgetTracker/backend/internal/service"
	"encoding/json"
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
	exchange := r.URL.Query().Get("exchange")

	var apiKeys []model.APIKey
	var err error

	if exchange != "" {
		key, err := h.service.GetAPIKey(ctx, exchange)
		if err != nil {
			http.Error(w, "API key not found", http.StatusNotFound)
			return
		}
		apiKeys = []model.APIKey{*key}
	} else {
		apiKeys, err = h.service.GetAllAPIKeys(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

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

	for _, key := range keys {
		if err := h.service.SaveAPIKey(ctx, &key); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
