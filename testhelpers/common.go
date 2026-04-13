package testhelpers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/selfshop-dev/ms-catalog/internal/handler"
)

func TimeNow() time.Time { return time.Now().UTC().Truncate(time.Millisecond) }
func UUIDNew() uuid.UUID { return uuid.New() }

// NewJSONRequest creates an HTTP request with a JSON-encoded body.
// Uses NewRequestWithContext with context.Background() as the base context.
func NewJSONRequest(t *testing.T, method, path string, body any) *http.Request {
	t.Helper()
	b, err := json.Marshal(body)
	require.NoError(t, err)
	r := httptest.NewRequestWithContext(
		context.Background(),
		method, path,
		bytes.NewReader(b),
	)
	r.Header.Set("Content-Type", "application/json")
	return r
}

// ServeWithChi registers the route on a chi router and serves the request.
// Required because chi.URLParam only resolves path parameters when the request
// passes through a chi router — direct handler calls leave the params empty.
func ServeWithChi(t *testing.T, route handler.Route, r *http.Request) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	router := chi.NewRouter()
	router.Method(route.Method, route.Path, route.Handler)
	router.ServeHTTP(w, r)
	return w
}
