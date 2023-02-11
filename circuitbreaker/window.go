package circuitbreaker

type WindowCallOutcome int

const (
	WindowCallSuccess WindowCallOutcome = iota
	WindowCallFailed
)

type Window interface {
	SuccessRate() float64
	Notify(WindowCallOutcome)
}

type IncrementalCountBasedWindow struct {
	samples []WindowCallOutcome
	head    int
	success int
	failure int
}

func NewCountBasedWindow(samples int) *IncrementalCountBasedWindow {
	outcomes := make([]WindowCallOutcome, samples)
	for i := range outcomes {
		outcomes[i] = -1
	}
	return &IncrementalCountBasedWindow{
		samples: outcomes,
	}
}

func (c *IncrementalCountBasedWindow) SuccessRate() float64 {
	return float64(c.success) / float64(c.success+c.failure)
}

func (c *IncrementalCountBasedWindow) Notify(outcome WindowCallOutcome) {
	var curr = c.samples[c.head]
	if curr == WindowCallSuccess {
		c.success--
	} else if curr == WindowCallFailed {
		c.failure--
	}

	if outcome == WindowCallSuccess {
		c.success++
	} else {
		c.failure++
	}

	c.samples[c.head] = outcome
	c.head = (c.head + 1) % len(c.samples)
}
