package notifications
import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/codestar-sunil/push-notif/internal/httpx"
)


type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler{
	return &Handler{repo: repo}
}

func (h *Handler) Routes(r chi.Router) {
	r.Post("/notifications", h.CreateNotification)
	r.Get("/notifications/{id}", h.GetNotification)
}


func (h *Handler) CreateNotification(w http.ResponseWriter, r *http.Request) {
	var req CreateNotificationRequest
	if err:=json.NewDecoder(r.Body).Decode(&req); err!=nil{
		httpx.WriteError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err:=req.Validate(); err!=nil{
		httpx.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	notification, err := h.repo.CreateNotification(r.Context(), req.ToNotification())
	if err != nil {
		httpx.WriteError(w, "Failed to create notification", http.StatusInternalServerError)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, notification)
}


func (h *Handler) GetNotification(w http.ResponseWriter, r *http.Request) {
	idParam:=chi.URLParam(r,"id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		httpx.WriteError(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}
	notification,err:=h.repo.GetNotificationByID(r.Context(),id)
	if err != nil {
		httpx.WriteError(w, "Notification not found", http.StatusNotFound)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, notification)
}