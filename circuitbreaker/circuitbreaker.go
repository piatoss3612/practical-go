package circuitbreaker

import (
	"context"
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

type Circuit func(ctx context.Context) (interface{}, error)

func (cb *CircuitBreaker) Exec(ctx context.Context, c Circuit) (interface{}, error) {
	panic("implement me")
}

func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.counter.reset()
	_ = cb.setState(StateClosed)
}

func (cb *CircuitBreaker) State() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

func (cb *CircuitBreaker) SetState(state State) bool {
	if !state.Valid() {
		return false
	}

	return cb.setStateWithLock(state)
}

func (cb *CircuitBreaker) setStateWithLock(state State) bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.setState(state)
}

func (cb *CircuitBreaker) setState(state State) bool {
	if cb.state == state {
		return false
	}

	cb.state = state
	if cb.onStateChange != nil {
		cb.onStateChange(cb.state, state)
	}

	return true
}

func (cb *CircuitBreaker) Counter() Counter {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.counter
}
