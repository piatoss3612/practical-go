package circuitbreaker

import "time"

var (
	DefaultHalfOpenMaxSuccesses = uint32(5)
	DefaultClearInterval        = 1 * time.Second
	DefaultOpenTimeout          = 60 * time.Second
	DefaultTrip                 = func(c Counter) bool {
		return c.TotalFailures > 5
	}
	DefaultOnStateChange = func(from, to State) {}
)

type Option func(c *CircuitBreaker) // CircuitBreaker의 옵션을 설정하는 함수의 타입

// halfOpenMaxSuccesses를 설정하는 옵션 함수
func WithHalfOpenMaxSuccesses(n uint32) Option {
	return func(c *CircuitBreaker) {
		c.halfOpenMaxSuccesses = n
	}
}

// clearInterval을 설정하는 옵션 함수
func WithClearInterval(d time.Duration) Option {
	return func(c *CircuitBreaker) {
		c.clearInterval = d
	}
}

// openTimeout을 설정하는 옵션 함수
func WithOpenTimeout(d time.Duration) Option {
	return func(c *CircuitBreaker) {
		c.openTimeout = d
	}
}

// trip을 설정하는 옵션 함수
func WithTripFunc(f TripFunc) Option {
	return func(c *CircuitBreaker) {
		c.trip = f
	}
}

// onStateChange를 설정하는 옵션 함수
func WithStateChangeHook(f StateChangeHook) Option {
	return func(c *CircuitBreaker) {
		c.onStateChange = f
	}
}

// 기본 옵션을 설정하는 옵션 함수
func WithDefaultOptions() Option {
	return func(c *CircuitBreaker) {
		c.halfOpenMaxSuccesses = DefaultHalfOpenMaxSuccesses
		c.clearInterval = DefaultClearInterval
		c.openTimeout = DefaultOpenTimeout
		c.trip = DefaultTrip
		c.onStateChange = DefaultOnStateChange
	}
}
