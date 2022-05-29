package router

import (
	"net/http"
)

func (h *Handler) SleepDebug(rw http.ResponseWriter, r *http.Request) {
	data, err := h.App.SleepDebug()
	JsonCheckError(rw, r, http.StatusOK, data, err)
}
