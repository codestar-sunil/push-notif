package devices

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/codestar-sunil/push-notif/internal/httpx"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Routes(r chi.Router) {
	r.Post("/devices/register", h.Register)
	r.Patch("/devices/{id}/heartbeat", h.Heartbeat)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		httpx.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	device, err := h.repo.Upsert(r.Context(), req.ToDevice())
	if err != nil {
		httpx.WriteError(w, "failed to register device", http.StatusInternalServerError)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, device)
}

func (h *Handler) Heartbeat(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		httpx.WriteError(w, "Invalid device ID", http.StatusBadRequest)
		return
	}

	err = h.repo.Heartbeat(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrDeviceNotFound) {
			httpx.WriteError(w, "Device not found", http.StatusNotFound)
			return
		}
		httpx.WriteError(w, "Failed to update heartbeat", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
