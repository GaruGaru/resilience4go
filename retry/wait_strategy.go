package retry

import "time"

// WaitStrategy provide a wait strategy for the retry policy controlling the delay
type WaitStrategy interface {
	OnError()
	Reset()
	Delay() time.Duration
}

type FixedDelay struct {
	delay time.Duration
}

func NewFixedDelay(delay time.Duration) *FixedDelay {
	return &FixedDelay{delay: delay}
}

func (f *FixedDelay) OnError() {
}

func (f *FixedDelay) Reset() {
}

func (f *FixedDelay) Delay() time.Duration {
	return f.delay
}

type ExponentialBackoff struct {
	base    time.Duration
	max     time.Duration
	current time.Duration
}

func NewExponentialBackoff(base time.Duration, max time.Duration) *ExponentialBackoff {
	return &ExponentialBackoff{base: base, max: max, current: base}
}

func (e *ExponentialBackoff) OnError() {
	if e.current == e.max {
		return
	}
	e.current = e.current * 2
	if e.current > e.max {
		e.current = e.max
	}
}

func (e *ExponentialBackoff) Reset() {
	e.current = e.base
}

func (e *ExponentialBackoff) Delay() time.Duration {
	return e.current
}
