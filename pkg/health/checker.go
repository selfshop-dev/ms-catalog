package health

import "context"

// Checker is the interface every health dependency must implement.
//
// Name identifies the dependency in the readyz JSON response and in logs.
// Prefer distinct names — duplicates make readiness responses and log output
// ambiguous.
//
// Check performs the actual health probe. It must honour ctx cancellation —
// the context carries the shared [Collector] timeout and the parent request
// context. Check must be safe for concurrent calls.
type Checker interface {
	Name() string
	Check(ctx context.Context) error
}
