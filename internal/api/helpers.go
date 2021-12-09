package api

//
// import (
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"strings"
// )
//
// type envelope map[string]interface{}
//
// func (a *main.application) readJSON(rw http.ResponseWriter, r *http.Request,
// 	data interface{}) error {
//
// 	maxBytes := 1_048_576 // 1MB
// 	r.Body = http.MaxBytesReader(rw, r.Body, int64(maxBytes))
//
// 	d := json.NewDecoder(r.Body)
// 	// d.DisallowUnknownFields()
//
// 	err := d.Decode(data)
// 	if err != nil {
// 		var syntaxError *json.SyntaxError
// 		var unmarshalTypeError *json.UnmarshalTypeError
// 		var invalidUnmarshalError *json.InvalidUnmarshalError
//
// 		switch {
// 		case errors.As(err, &syntaxError):
// 			return fmt.Errorf("body contains badly-formed JSON (at character %d)",
// 				syntaxError.Offset)
//
// 		// In some circumstances Decode() may also return an io.ErrUnexpectedEOF error
// 		// for syntax errors in the JSON. There is an open issue regarding this at
// 		// https://github.com/golang/go/issues/25956.
// 		case errors.Is(err, io.ErrUnexpectedEOF):
// 			return errors.New("body contains badly-formed JSON")
//
// 		case errors.As(err, &unmarshalTypeError):
// 			if unmarshalTypeError.Field != "" {
// 				return fmt.Errorf("body contains incorrect JSON type for field %q",
// 					unmarshalTypeError.Field)
// 			}
// 			return fmt.Errorf("body contains incorrect JSON type (at character %d)",
// 				unmarshalTypeError.Offset)
//
// 		case errors.Is(err, io.EOF):
// 			return errors.New("body must not be empty")
//
// 		// there's an open issue at https://github.com/golang/go/issues/29035
// 		// regarding turning this into a distinct error type in the future.
// 		case strings.HasPrefix(err.Error(), "json: unknown field "):
// 			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
// 			return fmt.Errorf("body contains unknown key %s", fieldName)
//
// 		// there is an open issue about turning
// 		// this into a distinct error type at https://github.com/golang/go/issues/30715.
// 		case err.Error() == "http: request body too large":
// 			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)
//
// 		case errors.As(err, &invalidUnmarshalError):
// 			panic(err)
//
// 		default:
// 			return err
// 		}
// 	}
//
// 	err = d.Decode(&struct{}{})
// 	if err != io.EOF {
// 		return errors.New("the body cannot have more than one JSON value")
// 	}
//
// 	return nil
// }
//
// func (a *main.application) writeJSON(rw http.ResponseWriter, status int, data envelope,
// 	headers http.Header) error {
// 	js, err := json.Marshal(data)
// 	if err != nil {
// 		return err
// 	}
//
// 	for k, v := range headers {
// 		rw.Header()[k] = v
// 	}
//
// 	rw.Header().Set("Content-Type", "application/json")
// 	rw.WriteHeader(status)
// 	rw.Write(js)
//
// 	return nil
// }
//
// // The background() helper accepts an arbitrary function as a parameter.
// func (a *main.application) background(fn func()) {
// 	// Launch a background goroutine.
// 	go func() {
// 		// Recover any panic.
// 		defer func() {
// 			if err := recover(); err != nil {
// 				//a.logger.PrintError(fmt.Errorf("%s", err), nil)
// 				a.logger.Println(err)
// 			}
// 		}()
//
// 		// Execute the arbitrary function that we passed as the parameter.
// 		fn()
// 	}()
// }
//
