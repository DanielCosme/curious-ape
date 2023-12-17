package transport

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"gotest.tools/v3/assert"
)

func TestPing(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NilError(t, err)

	ping(rr, r)
	rs := rr.Result()

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)

	assert.NilError(t, err)
	assert.Equal(t, rs.StatusCode, http.StatusOK)
	assert.Equal(t, string(body), "OK")
}
