package devices

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	device, err := h.repo.Upsert(r.Context(), req.ToDevice())
	if err != nil {
		writeError(w, "failed to register device", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, device)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, message string, statusCode int) {
	writeJSON(w, statusCode, map[string]string{"error": message})
}
