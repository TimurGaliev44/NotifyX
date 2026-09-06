package handler

import (
	"encoding/json"
	"net/http"

	"github.com/TimurGaliev44/notifyx/internal/domain"
	"github.com/TimurGaliev44/notifyx/internal/ingestion/service"
	"github.com/go-chi/chi/v5"
)

type createNotificationRequest struct {
	Channel        string         `json:"channel"`
	Recipient      string         `json:"recipient"`
	TemplateID     string         `json:"template_id"`
	Payload        map[string]any `json:"payload"`
	IdempotencyKey string         `json:"idempotency_key"`
}

type Server struct {
	router  *chi.Mux
	service *service.Service
}

func New(svc *service.Service) *Server {
	s := &Server{router: chi.NewRouter(), service: svc}
	s.router.Post("/notifications", s.createNotification)
	return s
}

func (s *Server) createNotification(w http.ResponseWriter, r *http.Request) {
	var req createNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Channel != string(domain.ChannelEmail) && req.Channel != string(domain.ChannelSMS) {
		http.Error(w, "invalid channel", http.StatusBadRequest)
		return
	}
	if req.Recipient == "" || req.TemplateID == "" {
		http.Error(w, "recipient and template_id are required", http.StatusBadRequest)
		return
	}

	n := &domain.Notification{
		Channel:        domain.Channel(req.Channel),
		Recipient:      req.Recipient,
		TemplateID:     req.TemplateID,
		Payload:        req.Payload,
		IdempotencyKey: req.IdempotencyKey,
		Source:         "ingestion-api",
	}

	created, err := s.service.Create(r.Context(), n)
	if err != nil {
		http.Error(w, "failed to create notification", http.StatusInternalServerError)
		return
	}
	status := http.StatusOK
	if created {
		status = http.StatusCreated
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"id": n.ID.String()})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
