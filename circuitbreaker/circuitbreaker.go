package circuitbreaker

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"
)

var (
	ErrOpenState    = errors.New("circuit breaker is in open state")
	ErrInvalidState = errors.New("circuit breaker is in invalid state")
)

type TripFunc func(c Counter) bool

type StateChangeHook func(from, to State)

type CircuitBreaker struct {
	halfOpenMaxSuccesses uint32          // max successes in half-open state
	clearInterval        time.Duration   // counter clear interval in closed state
	openTimeout          time.Duration   // timeout in open state
	trip                 TripFunc        // trip function to determine if the circuit should be tripped
	onStateChange        StateChangeHook // hook to be called when the state changes

	state   State        // current state
	counter Counter      // counter to track the number of requests
	mu      sync.RWMutex // mutex to protect the state and counter
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

	go cb.resetCounterInterval()

	return cb
}

type Circuit func(ctx context.Context) (interface{}, error)

func (cb *CircuitBreaker) Execute(ctx context.Context, c Circuit) (interface{}, error) {
	// check if the circuit breaker is ready
	if err := cb.ready(); err != nil {
		return nil, err
	}

	// execute the circuit
	res, err := c(ctx)

	// record the result
	return res, cb.done(err)
}

func (cb *CircuitBreaker) ready() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed, StateHalfOpen:
		return nil
	case StateOpen:
		return ErrOpenState
	default:
		return ErrInvalidState
	}
}

func (cb *CircuitBreaker) done(err error) error {
	if err == nil {
		cb.success()
		return nil
	}

	cb.fail()
	return err
}

func (cb *CircuitBreaker) success() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		cb.counter.onSuccess()
	case StateHalfOpen:
		cb.counter.onSuccess()
		if cb.counter.TotalSuccesses >= cb.halfOpenMaxSuccesses {
			_ = cb.setState(StateClosed)
			go cb.resetCounterInterval()
		}
	}
}

func (cb *CircuitBreaker) resetCounterInterval() {
	ticker := time.NewTicker(cb.clearInterval)

	for range ticker.C {
		cb.mu.RLock()
		if cb.state != StateClosed {
			cb.mu.RUnlock()
			ticker.Stop()
			return
		}
		cb.mu.RUnlock()

		cb.mu.Lock()
		cb.reset()
		cb.mu.Unlock()

		slog.Info("Successfully reset circuit breaker counter")
	}
}

func (cb *CircuitBreaker) fail() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		cb.counter.onFailure()
		if cb.trip(cb.counter) {
			_ = cb.setState(StateOpen)
			go cb.checkOpenTimeout()
		}
	case StateHalfOpen:
		_ = cb.setState(StateOpen)
		go cb.checkOpenTimeout()
	}
}

func (cb *CircuitBreaker) checkOpenTimeout() {
	<-time.NewTimer(cb.openTimeout).C

	cb.mu.Lock()
	_ = cb.setState(StateHalfOpen)
	cb.mu.Unlock()
}

func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.reset()
}

func (cb *CircuitBreaker) reset() {
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

func (cb *CircuitBreaker) setState(newState State) bool {
	current := cb.state

	if current == newState {
		return false
	}

	cb.state = newState
	if cb.onStateChange != nil {
		cb.onStateChange(current, newState)
	}

	return true
}

func (cb *CircuitBreaker) Counter() Counter {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.counter
}
