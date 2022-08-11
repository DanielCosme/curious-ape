package routes

import (
	"github.com/danielcosme/curious-ape/rest"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (h *Handler) Oauth2Connect(rw http.ResponseWriter, r *http.Request) {
	url, err := h.App.Oauth2ConnectProvider(chi.URLParam(r, "provider"))
	if err != nil {
		rest.ErrResponse(rw, http.StatusInternalServerError, err)
		return
	}

	headers := http.Header{}
	headers.Set("location", url)
	rest.JSONWithHeaders(rw, http.StatusTemporaryRedirect, nil, headers)
}

func (h *Handler) Oauth2Success(rw http.ResponseWriter, r *http.Request) {
	code := r.Form.Get("code")
	err := h.App.Oauth2Success(chi.URLParam(r, "provider"), code)
	if err != nil {
		rest.ErrInternalServer(rw)
		return
	}

	rest.JSONStatusOk(rw, envelope{"message": "ok"})
}

func (h *Handler) AddToken(rw http.ResponseWriter, r *http.Request) {
	msg, err := h.App.Oauth2AddAPIToken(r.Form.Get("token"), chi.URLParam(r, "provider"))
	JsonCheckError(rw, r, http.StatusCreated, envelope{"valid token for": msg}, err)
}
