package transport

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (h *Transport) Oauth2ConnectForm(w http.ResponseWriter, r *http.Request) {
	h.render(w, http.StatusOK, "auth.gohtml", h.newTemplateData(r))
}

func (h *Transport) Oauth2Connect(w http.ResponseWriter, r *http.Request) {
	url, err := h.App.Oauth2ConnectProvider(chi.URLParam(r, "provider"))
	if err != nil {
		h.serverError(w, err)
		return
	}
	h.App.Log.Debug("url to redirect: ", url)

	w.Header().Set("location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Transport) Oauth2Success(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.serverError(w, err)
	}

	code := r.Form.Get("code")
	prov := chi.URLParam(r, "provider")
	if code == "" || prov == "" {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	err := h.App.Oauth2Success(prov, code)
	if err != nil {
		h.serverError(w, err)
		return
	}

	h.App.Log.Info("successfully authenticated with:", prov)
	w.WriteHeader(http.StatusOK)
}

func (h *Transport) AddToken(rw http.ResponseWriter, r *http.Request) {
	msg, err := h.App.AuthAddAPIToken(r.Form.Get("token"), chi.URLParam(r, "provider"))
	if err != nil {
		h.serverError(rw, err)
		return
	}
	h.App.Log.Info("valid token for: ", msg)
	rw.WriteHeader(http.StatusCreated)
}
