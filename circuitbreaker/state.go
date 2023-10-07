package circuitbreaker

type State int // 서킷 브레이커의 상태를 나타내는 타입

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

// 서킷 브레이커의 상태를 문자열로 반환하는 메서드 (Stringer 인터페이스 구현)
func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// 서킷 브레이커의 상태가 유효한지 확인하는 메서드
func (s State) Valid() bool {
	return s == StateClosed || s == StateOpen || s == StateHalfOpen
}
