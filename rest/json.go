package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Envelope map[string]interface{}

const HeaderContentType = "Content-Type"

func JSON(rw http.ResponseWriter, status int, data Envelope) error {
	return writeJSON(rw, status, data, nil)
}

func JSONWithHeaders(rw http.ResponseWriter, status int, data Envelope, headers http.Header) error {
	return writeJSON(rw, status, data, headers)
}

func writeJSON(rw http.ResponseWriter, status int, data Envelope, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	WriteHeaders(rw, headers)
	rw.Header().Set(HeaderContentType, "application/json")
	rw.WriteHeader(status)
	_, err = rw.Write(js)
	return err
}

func WriteHeaders(rw http.ResponseWriter, headers http.Header) {
	for k, v := range headers {
		rw.Header()[k] = v
	}
}

func ReadJSON(r *http.Request, data interface{}) error {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	err := d.Decode(data)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		// In some circumstances Decode() may also return an io.ErrUnexpectedEOF error
		// for syntax errors in the JSON. There is an open issue regarding this at
		// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// there's an open issue at https://github.com/golang/go/issues/29035
		// regarding turning this into a distinct error type in the future.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = d.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("the body cannot have more than one JSON value")
	}

	return nil
}
