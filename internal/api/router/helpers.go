package router

import (
	"github.com/danielcosme/curious-ape/internal/core/database"
	"net/http"

	"github.com/danielcosme/curious-ape/rest"
	"github.com/danielcosme/curious-ape/sdk/errors"
)

func JsonCheckError(rw http.ResponseWriter, r *http.Request, status int, data interface{}, err error) {
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			rest.ErrResponse(rw, http.StatusNotFound, err)
			return
		} else {
			rest.ErrResponse(rw, http.StatusInternalServerError, err)
			return
		}
	}

	rest.JSON(rw, status, data)
}

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
