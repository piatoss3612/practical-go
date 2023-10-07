package circuitbreaker

import "time"

var (
	DefaultHalfOpenMaxSuccesses = uint32(5)
	DefaultOpenTimeout          = 60 * time.Second
	DefaultTrip                 = func(c Counter) bool {
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

func WithOpenTimeout(d time.Duration) Option {
	return func(c *CircuitBreaker) {
		c.openTimeout = d
	}
}

func WithTripFunc(f TripFunc) Option {
	return func(c *CircuitBreaker) {
		c.trip = f
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
		c.openTimeout = DefaultOpenTimeout
		c.trip = DefaultTrip
		c.onStateChange = DefaultOnStateChange
	}
}
