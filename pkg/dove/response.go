package dove

import (
	"net/http"
)

type Response struct {
	http.ResponseWriter
}

func (r *Response) Header() http.Header {
	return nil
}

func (r *Response) Write([]byte) (int, error) {
	return 0, nil
}

func (r *Response) WriteHeader(statusCode int) {

}
