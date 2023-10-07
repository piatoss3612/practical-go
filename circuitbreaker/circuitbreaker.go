package circuitbreaker

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrOpenState    = errors.New("circuit breaker is in open state")
	ErrInvalidState = errors.New("circuit breaker is in invalid state")
)

type TripFunc func(c Counter) bool // 서킷 브레이커가 open 상태로 전환되기 위한 조건을 판단하는 함수의 타입

type StateChangeHook func(from, to State) // 서킷 브레이커의 상태가 변경될 때 호출되는 함수의 타입

// 서킷 브레이커 구조체
type CircuitBreaker struct {
	halfOpenMaxSuccesses uint32          // half open 상태에서 closed 상태로 전환되기 위한 최소 성공 횟수
	openTimeout          time.Duration   // open 상태에서 half open 상태로 전환되기 위한 시간
	trip                 TripFunc        // 서킷 브레이커가 open 상태로 전환되기 위한 조건을 판단하는 함수
	onStateChange        StateChangeHook // 서킷 브레이커의 상태가 변경될 때 호출되는 함수

	state   State        // 서킷 브레이커의 상태
	counter Counter      // 서킷 브레이커의 상태를 판단하기 위한 counter (성공/실패 횟수)
	mu      sync.RWMutex // 서킷 브레이커의 상태를 변경하기 위한 mutex
}

// CircuitBreaker의 생성자 함수.
// opts에는 CircuitBreaker의 옵션을 설정하는 함수들이 들어간다.
func New(opts ...Option) *CircuitBreaker {
	cb := &CircuitBreaker{
		state:   StateClosed,    // 서킷 브레이커의 초기 상태는 closed 상태이다.
		counter: Counter{},      // 서킷 브레이커의 counter를 초기화한다.
		mu:      sync.RWMutex{}, // 서킷 브레이커의 상태를 변경하기 위한 mutex를 초기화한다.
	}

	WithDefaultOptions()(cb) // 기본 옵션을 설정한다.

	// opts에 있는 옵션들을 설정한다.
	for _, opt := range opts {
		opt(cb)
	}

	return cb
}

type Circuit func(ctx context.Context) (interface{}, error) // 서킷 브레이커로 wrapping될 함수의 타입

// 서킷 브레이커로 wrapping된 함수를 실행하는 메서드.
// ctx, c는 각각 wrapping될 함수로 전달되는 context와 wrapping될 함수이다.
func (cb *CircuitBreaker) Execute(ctx context.Context, c Circuit) (interface{}, error) {
	// 서킷 브레이커가 c를 실행할 수 있는 상태인지 확인
	if err := cb.ready(); err != nil {
		return nil, err
	}

	// c를 실행
	res, err := c(ctx)

	// c의 실행 결과에 따라 서킷 브레이커의 상태를 변경
	return res, cb.done(err)
}

// 서킷 브레이커가 c를 실행할 수 있는 상태인지 확인하는 메서드
func (cb *CircuitBreaker) ready() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed, StateHalfOpen: // 서킷 브레이커가 closed 상태이거나 half open 상태인 경우
		return nil
	case StateOpen: // 서킷 브레이커가 open 상태인 경우
		return ErrOpenState
	default: // 서킷 브레이커가 알 수 없는 상태인 경우
		return ErrInvalidState
	}
}

// c의 실행 결과에 따라 서킷 브레이커의 상태를 변경하는 메서드
func (cb *CircuitBreaker) done(err error) error {
	// err가 nil이라면 c의 실행이 성공
	if err == nil {
		cb.success()
		return nil
	}

	// err가 nil이 아니라면 c의 실행이 실패
	cb.fail()
	return err
}

// 서킷 브레이커가 c를 실행한 결과가 성공인 경우 호출되는 메서드
func (cb *CircuitBreaker) success() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed: // 서킷 브레이커가 closed 상태인 경우
		cb.counter.onSuccess() // 성공 횟수를 증가시킨다.
	case StateHalfOpen: // 서킷 브레이커가 half open 상태인 경우
		cb.counter.onSuccess() // 성공 횟수를 증가시킨다.

		// half open 상태에서 closed 상태로 전환되기 위한 최소 성공 횟수를 만족하는지 확인한다.
		if cb.counter.TotalSuccesses >= cb.halfOpenMaxSuccesses {
			_ = cb.setState(StateClosed) // closed 상태로 전환한다.
			cb.counter.resetFailures()   // 실패 횟수를 초기화한다.
		}
	}
}

// 서킷 브레이커가 c를 실행한 결과가 실패인 경우 호출되는 메서드
func (cb *CircuitBreaker) fail() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed: // 서킷 브레이커가 closed 상태인 경우
		cb.counter.onFailure() // 실패 횟수를 증가시킨다.

		// 서킷 브레이커가 open 상태로 전환되기 위한 조건을 만족하는지 확인한다.
		if cb.trip(cb.counter) {
			_ = cb.setState(StateOpen) // open 상태로 전환한다.
			go cb.checkOpenTimeout()   // open 상태에서 half open 상태로 전환되기 위한 시간을 체크하는 goroutine을 실행한다.
		}
	case StateHalfOpen: // 서킷 브레이커가 half open 상태인 경우
		_ = cb.setState(StateOpen) // open 상태로 전환한다.
		go cb.checkOpenTimeout()   // open 상태에서 half open 상태로 전환되기 위한 시간을 체크하는 goroutine을 실행한다.
	}
}

// open 상태에서 half open 상태로 전환되기 위한 시간을 체크하는 goroutine
func (cb *CircuitBreaker) checkOpenTimeout() {
	<-time.NewTimer(cb.openTimeout).C // open 상태에서 half open 상태로 전환되기 위한 시간이 지날 때까지 대기한다.

	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.setState(StateHalfOpen)  // half open 상태로 전환한다.
	cb.counter.resetSuccesses() // 성공 횟수를 초기화한다.
}

// 서킷 브레이커의 상태를 초기화하는 메서드
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.counter.reset()           // 카운터를 초기화한다.
	_ = cb.setState(StateClosed) // 서킷 브레이커의 상태를 closed 상태로 전환한다.
}

// 서킷 브레이커의 상태를 반환하는 메서드
func (cb *CircuitBreaker) State() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// 서킷 브레이커의 상태를 설정하는 메서드
func (cb *CircuitBreaker) SetState(state State) bool {
	// state가 유효한지 확인한다.
	if !state.Valid() {
		return false
	}

	return cb.setStateWithLock(state) // 서킷 브레이커의 상태를 설정한다.
}

// 서킷 브레이커의 상태를 설정하는 메서드 (lock을 사용)
func (cb *CircuitBreaker) setStateWithLock(state State) bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.setState(state)
}

// 서킷 브레이커의 상태를 설정하는 메서드 (lock을 사용하지 않음)
func (cb *CircuitBreaker) setState(newState State) bool {
	current := cb.state // 현재 상태를 저장한다.

	// 현재 상태와 새로운 상태가 같다면 false를 반환한다.
	if current == newState {
		return false
	}

	cb.state = newState // 서킷 브레이커의 상태를 새로운 상태로 변경한다.
	// 서킷 브레이커의 상태가 변경될 때 호출되는 함수가 있다면 호출한다.
	if cb.onStateChange != nil {
		cb.onStateChange(current, newState)
	}

	return true
}

// 서킷 브레이커의 counter를 반환하는 메서드
func (cb *CircuitBreaker) Counter() Counter {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.counter
}
