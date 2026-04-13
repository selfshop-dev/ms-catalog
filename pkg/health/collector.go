package health

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
)

// Collector runs a set of [Checker]s concurrently and aggregates results.
// Construct one with [New]; the zero value is not usable.
type Collector struct {
	logr     *zap.Logger
	checkers []Checker
	conf     Config
}

// New creates a [Collector] with the provided config, logger, and checkers.
//
// Returns an error if any element of checkers is nil — a nil [Checker] cannot
// provide a name and would panic inside [Collector.CheckAll] rather than
// producing a recoverable error result.
func New(c Config, l *zap.Logger, cs ...Checker) (*Collector, error) {
	for i, c := range cs {
		if c == nil {
			return nil, fmt.Errorf("health: checkers[%d] is nil", i)
		}
	}

	return &Collector{
		conf:     c,
		logr:     l.Named("health"),
		checkers: cs,
	}, nil
}

// CheckAll runs all registered checkers concurrently and returns aggregated
// results. It is safe to call from multiple goroutines simultaneously.
//
// Execution model:
//   - A context with [Config.Timeout] is derived from ctx and shared by all
//     checker goroutines. Callers should still pass a meaningful parent context
//     (e.g. the request context) so that upstream cancellation propagates.
//   - Up to [Config.MaxConcurrency] checkers run at a time; the rest wait for
//     the semaphore. If the context expires while waiting, the remaining
//     checkers are marked as "not started: <reason>" and HasErrors is set.
//   - Results are written to pre-allocated index positions — no mutex needed
//     for the slice itself. HasErrors is updated atomically.
//   - Panics inside checkers are recovered by [runCheck] and reported as errors.
func (c *Collector) CheckAll(ctx context.Context) CheckAllResult {
	n := len(c.checkers)
	if n == 0 {
		return CheckAllResult{}
	}

	ctx, cancel := context.WithTimeout(ctx, c.conf.Timeout)
	defer cancel()

	// Pre-fill every result with the checker name and a "not started" placeholder.
	// Checkers that are never started keep this; it is replaced with the real
	// context error at the point of semaphore failure so the response explains
	// exactly why each checker was skipped.
	rs := make([]CheckResult, n)
	for i, ch := range c.checkers {
		rs[i] = CheckResult{Name: ch.Name(), Error: "not started"}
	}

	var hasErrors atomic.Bool

	sem := semaphore.NewWeighted(c.conf.MaxConcurrency)
	var wg sync.WaitGroup

	for i := range c.checkers {
		// sem.Acquire blocks until a slot is available or ctx is done.
		// On failure, propagate the real context error to all remaining entries.
		if err := sem.Acquire(ctx, 1); err != nil {
			reason := "not started: " + ctx.Err().Error()
			for j := i; j < n; j++ {
				rs[j].Error = reason
			}
			hasErrors.Store(true)
			break
		}

		wg.Add(1)
		go func(i int) {
			defer sem.Release(1)
			defer wg.Done()

			r := runCheck(ctx, c.checkers[i], c.logr)
			// Safe without a mutex: each goroutine writes only to index i.
			// Visibility is guaranteed by wg.Wait() below.
			rs[i] = r
			if !r.ok() {
				hasErrors.Store(true)
			}
		}(i)
	}

	wg.Wait()

	return CheckAllResult{
		Results:   rs,
		HasErrors: hasErrors.Load(),
	}
}

// runCheck executes a single checker and captures panics.
// A panic is converted into a CheckResult with Error set to "<n>: panic".
// The panic value and stack trace are logged at Error level so the cause is
// visible in structured logs without leaking internal details via the public
// readiness HTTP response.
func runCheck(ctx context.Context, c Checker, l *zap.Logger) (r CheckResult) {
	name := c.Name()

	defer func() {
		if rec := recover(); rec != nil {
			l.Error("health checker panic",
				zap.String("checker", name),
				zap.Any("panic", rec),
				zap.Stack("stack"),
			)
			r = CheckResult{
				Name:  name,
				Error: name + ": panic",
			}
		}
	}()

	if err := c.Check(ctx); err != nil {
		return CheckResult{Name: name, Error: err.Error()}
	}
	return CheckResult{Name: name}
}
