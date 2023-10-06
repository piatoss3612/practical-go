package circuitbreaker

import "time"

var (
	DefaultHalfOpenMaxSuccesses = uint32(5)
	DefaultClearInterval        = 1 * time.Second
	DefaultOpenTimeout          = 60 * time.Second
	DefaultShouldTrip           = func(c Counter) bool {
		return c.TotalFailures > 5
	}
	DefaultOnStateChange = func(from, to State) {}
)

type Option func(c *CircuitBreaker)

func WithHalfOpenMaxSuccesses(n uint32) Option {
	return func(c *CircuitBreaker) {
		c.halfOpenMaxSuccesses = n
	}
}

func WithClearInterval(d time.Duration) Option {
	return func(c *CircuitBreaker) {
		c.clearInterval = d
	}
}

func WithOpenTimeout(d time.Duration) Option {
	return func(c *CircuitBreaker) {
		c.openTimeout = d
	}
}

func WithTripFunc(f TripFunc) Option {
	return func(c *CircuitBreaker) {
		c.shouldTrip = f
	}
}

func WithStateChangeHook(f StateChangeHook) Option {
	return func(c *CircuitBreaker) {
		c.onStateChange = f
	}
}

func WithDefaultOptions() Option {
	return func(c *CircuitBreaker) {
		c.halfOpenMaxSuccesses = DefaultHalfOpenMaxSuccesses
		c.clearInterval = DefaultClearInterval
		c.openTimeout = DefaultOpenTimeout
		c.shouldTrip = DefaultShouldTrip
		c.onStateChange = DefaultOnStateChange
	}
}
