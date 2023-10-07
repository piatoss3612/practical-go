package circuitbreaker

type Counter struct {
	TotalSuccesses       uint32
	TotalFailures        uint32
	ConsecutiveSuccesses uint32
	ConsecutiveFailures  uint32
}

func (c *Counter) onSuccess() {
	c.TotalSuccesses++
	c.ConsecutiveSuccesses++
	c.ConsecutiveFailures = 0
}

func (c *Counter) onFailure() {
	c.TotalFailures++
	c.ConsecutiveFailures++
	c.ConsecutiveSuccesses = 0
}

func (c *Counter) reset() {
	c.resetSuccesses()
	c.resetFailures()
}

func (c *Counter) resetSuccesses() {
	c.TotalSuccesses = 0
	c.ConsecutiveSuccesses = 0
}

func (c *Counter) resetFailures() {
	c.TotalFailures = 0
	c.ConsecutiveFailures = 0
}
