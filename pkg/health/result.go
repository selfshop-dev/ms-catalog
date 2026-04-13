package health

// CheckResult holds the outcome of a single [Checker] run.
// It is included in the JSON body of the readiness response.
type CheckResult struct {
	// Name is the value returned by [Checker.Name].
	Name string `json:"name"`

	// Error is the error message returned by [Checker.Check], or empty if the
	// check passed. Omitted from JSON when empty.
	Error string `json:"error,omitempty"`
}

// ok reports whether the check passed (no error).
func (r CheckResult) ok() bool { return r.Error == "" }

// CheckAllResult is the aggregated output of [Collector.CheckAll].
type CheckAllResult struct {
	// Results contains one entry per registered checker, in registration order.
	// Checkers that were not started because the shared context expired before
	// the semaphore could be acquired are reported with
	// Error = "not started: context deadline exceeded", so every entry always
	// carries a meaningful Name.
	Results []CheckResult

	// HasErrors is true if any checker failed or if the shared context expired
	// before all checkers could be started.
	HasErrors bool
}
