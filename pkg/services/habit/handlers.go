package habit

import "net/http"

type Handler struct {
	svc *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) Habits(w http.ResponseWriter, r *http.Request) {

}
