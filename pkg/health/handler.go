package health

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

var aliveBody = func() []byte {
	b, err := json.Marshal(struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	})
	if err != nil {
		panic("health: failed to encode alive body: " + err.Error())
	}
	return b
}()

// AliveHandler handles GET /alive.
// Returns 200 OK unconditionally — its only purpose is to prove the process
// is alive and the HTTP server is accepting connections. No external calls.
func (c *Collector) AliveHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(aliveBody) //nolint:errcheck // write error after WriteHeader is unrecoverable
}

// ReadyHandler handles GET /ready.
// Runs all registered [Checker]s concurrently via [Collector.CheckAll].
// Returns 200 when all pass, 503 Service Unavailable when any fail.
//
// The full check list is always included in the response body so operators
// can see the complete dependency status — not just the failures.
// Failed checks are logged at Warn level with structured fields.
func (c *Collector) ReadyHandler(w http.ResponseWriter, r *http.Request) {
	res := c.CheckAll(r.Context())

	status := http.StatusOK
	statusText := "ok"

	if res.HasErrors {
		status = http.StatusServiceUnavailable
		statusText = "unhealthy"
		c.logFailures(res.Results)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(struct { //nolint:errcheck // write error after WriteHeader is unrecoverable
		Status   string        `json:"status"`
		Checkers []CheckResult `json:"checkers"`
	}{
		Status:   statusText,
		Checkers: res.Results,
	})
}

func (c *Collector) logFailures(rs []CheckResult) {
	fs := make([]zap.Field, 0, len(rs))
	for _, r := range rs {
		if !r.ok() {
			fs = append(fs, zap.String(r.Name, r.Error))
		}
	}
	c.logr.Warn("readiness check failed", fs...)
}
