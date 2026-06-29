package segments

import (
	"encoding/json"
	"net/http"

	"github.com/codestar-sunil/push-notif/internal/httpx"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Routes(r chi.Router) {
	r.Post("/segments", h.Create)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateSegmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := req.Validate(); err != nil {
		httpx.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	segment, err := h.repo.CreateSegment(r.Context(), req.ToSegment())
	if err != nil {
		httpx.WriteError(w, "Failed to create segment", http.StatusInternalServerError)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, segment)
}
