package transport

//
// import (
// 	"github.com/danielcosme/curious-ape/cmd/api"
// 	"net/http"
// )
//
// // Show api information
// func (a *main.application) healthcheckerHandler(rw http.ResponseWriter, r *http.Request) {
// 	data := envelope{
// 		"status": "available",
// 		"systemInfo": map[string]string{
// 			"environment": "",
// 			"version":     main.version,
// 		},
// 	}
//
// 	err := a.writeJSON(rw, http.StatusOK, data, nil)
// 	if err != nil {
// 		a.serverErrorResponse(rw, r, err)
// 	}
// }
//
