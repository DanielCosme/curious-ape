package transport

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
)

type Status string

const (
	StatusSuccess Status = "success"
	StatusFail    Status = "fail"
	StatusError   Status = "error"
	StatusUnknown Status = "unknown"
)

type Payload struct {
	Status Status   `json:"status"`
	Data   envelope `json:"data"`
}

type envelope map[string]any

type EnvelopeError struct {
	Message string `json:"message"`
	Details any    `json:"details,omitzero"`
}

func JSONOK(w http.ResponseWriter, data envelope, headers http.Header) {
	JSON(w, http.StatusOK, data, headers)
}

func JSONError(w http.ResponseWriter, err error, status int) {
	slog.Error(err.Error())
	errPayload := &EnvelopeError{Message: err.Error()}

	// NOTE: This will fail for wrapped errors.
	switch e := err.(type) {
	case *ValidationError:
		errPayload.Details = e.v.details
	}
	JSON(w, status, envelope{"error": errPayload}, nil)
}

func JSON(w http.ResponseWriter, status int, data envelope, headers http.Header) {
	payload := Payload{
		Status: StatusUnknown,
		Data:   data,
	}
	switch status {
	case http.StatusOK, http.StatusCreated:
		payload.Status = StatusSuccess
	case http.StatusBadRequest:
		payload.Status = StatusFail
	case http.StatusInternalServerError:
		payload.Status = StatusError
	}

	js, err := json.Marshal(payload)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for k, v := range headers {
		w.Header()[k] = v
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
}

func Bind(body io.ReadCloser, out any) error {
	raw, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	defer body.Close()
	if len(raw) == 0 {
		return errors.New("empty body")
	}
	if !json.Valid(raw) {
		return errors.New("invalid json")
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return err
	}
	if validation, ok := out.(Validation); ok {
		v := NewValidator()
		validation.Validate(v)
		if !v.Valid() {
			return &ValidationError{v}
		}
	}
	return nil
}
