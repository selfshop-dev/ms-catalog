package handler

import (
	"net/http"

	response "github.com/selfshop-dev/lib-response"

	"github.com/selfshop-dev/ms-catalog/pkg/server/middleware"
)

// Respond is the package-level HTTP writer used by all handlers.
var Respond = response.NewWriter(func(r *http.Request) map[string]any {
	if id := middleware.RequestIDFromContext(r.Context()); id != "" {
		return map[string]any{"request_id": id}
	}
	return nil
})

// Not covered by a dedicated unit test because:
//   - response.NewWriter and middleware.RequestIDFromContext are third-party / shared
//     code with their own test suites — testing that we call them correctly is
//     testing configuration, not behaviour.
//   - The meta injection logic (request_id present/absent) is exercised indirectly
//     by every handler test that goes through Respond.
