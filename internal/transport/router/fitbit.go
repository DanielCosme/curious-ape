package router

import (
	"github.com/danielcosme/curious-ape/rest"
	"net/http"
)

func (h *Handler) Oauth2FitbitConnect(rw http.ResponseWriter, r *http.Request) {
	url, err := h.App.Oauth2ConnectProvider("fitbit")
	if err != nil {
		rest.ErrResponse(rw, http.StatusInternalServerError, err)
		return
	}

	headers := http.Header{}
	headers.Set("location", url)
	rest.JSONWithHeaders(rw, http.StatusTemporaryRedirect, nil, headers)
}

func (h *Handler) Oauth2FitbitSuccess(rw http.ResponseWriter, r *http.Request) {
	code := r.Form.Get("code")
	err := h.App.Oauth2Success("fitbit", code)
	if err != nil {
		rest.ErrInternalServer(rw)
		return
	}

	rest.JSONStatusOk(rw, nil)
}

func (h *Handler) SleepDebug(rw http.ResponseWriter, r *http.Request) {
	data, err := h.App.SleepDebug()
	JsonCheckError(rw, r, http.StatusOK, data, err)
}
