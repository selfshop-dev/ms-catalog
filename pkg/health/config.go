package health

import "time"

// Config holds tunable parameters for a [Collector].
// Both fields have safe defaults applied by [New] when zero.
type Config struct {
	// MaxConcurrency is the maximum number of [Checker]s that run in parallel.
	// Default: 8.
	MaxConcurrency int64

	// Timeout is the maximum wall-clock time for a full [Collector.CheckAll]
	// call, shared across all checkers.
	// Default: 6s.
	Timeout time.Duration
}

func DefaultConfig() Config {
	return Config{
		MaxConcurrency: 8,
		Timeout:        6 * time.Second,
	}
}
