package circuitbreaker

import (
	"sync"
	"time"
)

type TripFunc func(c Counter) bool

type StateChangeHook func(from, to State)

type CircuitBreaker struct {
	halfOpenMaxSuccesses uint32
	clearInterval        time.Duration
	openTimeout          time.Duration
	shouldTrip           TripFunc
	onStateChange        StateChangeHook

	state   State
	counter Counter
	mu      sync.RWMutex
}

func New(opts ...Option) *CircuitBreaker {
	cb := &CircuitBreaker{
		state:   StateClosed,
		counter: Counter{},
		mu:      sync.RWMutex{},
	}

	WithDefaultOptions()(cb)

	for _, opt := range opts {
		opt(cb)
	}

	return cb
}
