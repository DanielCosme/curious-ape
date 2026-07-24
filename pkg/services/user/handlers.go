package user

import (
	"log/slog"
	"net/http"

	"danicos.dev/daniel/curious-ape/pkg/services/user/pages"
	"danicos.dev/daniel/curious-ape/pkg/ui"
	"danicos.dev/daniel/curious-ape/pkg/web"
	"github.com/alexedwards/scs/v2"
)

type Handler struct {
	svc     *Service
	session *scs.SessionManager
}

func NewHandler(s *Service, session *scs.SessionManager) *Handler {
	return &Handler{svc: s, session: session}
}

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	if IsAuthenticated(r) {
		web.Redirect(w, "/")
		return
	}
	s := ui.NewUIState("Login")
	err := pages.LoginPage(s).Render(r.Context(), w)
	if err != nil {
		web.ErrInternalServer(err, w)
	}
}

func (h *Handler) LoginPost(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	id, err := h.svc.Authenticate(username, password)
	if err != nil {
		web.Err(http.StatusUnauthorized, err, w)
		return
	}

	err = h.session.RenewToken(r.Context())
	if err != nil {
		web.ErrInternalServer(err, w)
		return
	}

	slog.Info("User authenticated")
	h.session.Put(r.Context(), string(ctxKeyAuthenticatedUserID), id)
	web.Redirect(w, "/")
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.session.RenewToken(r.Context()); err != nil {
		web.ErrInternalServer(err, w)
		return
	}
	h.session.Remove(r.Context(), string(ctxKeyAuthenticatedUserID))
	web.Redirect(w, "/login")
}
