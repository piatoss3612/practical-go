package circuitbreaker

// 서킷 브레이커의 상태를 판단하기 위한 Counter (성공/실패 횟수) 구조체
type Counter struct {
	TotalSuccesses       uint32
	TotalFailures        uint32
	ConsecutiveSuccesses uint32
	ConsecutiveFailures  uint32
}

// 성공 횟수를 증가시키는 메서드
func (c *Counter) onSuccess() {
	c.TotalSuccesses++
	c.ConsecutiveSuccesses++
	c.ConsecutiveFailures = 0
}

// 실패 횟수를 증가시키는 메서드
func (c *Counter) onFailure() {
	c.TotalFailures++
	c.ConsecutiveFailures++
	c.ConsecutiveSuccesses = 0
}

// 카운터를 초기화하는 메서드
func (c *Counter) reset() {
	c.resetSuccesses()
	c.resetFailures()
}

// 성공 횟수를 초기화하는 메서드
func (c *Counter) resetSuccesses() {
	c.TotalSuccesses = 0
	c.ConsecutiveSuccesses = 0
}

// 실패 횟수를 초기화하는 메서드
func (c *Counter) resetFailures() {
	c.TotalFailures = 0
	c.ConsecutiveFailures = 0
}
