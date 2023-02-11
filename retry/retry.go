package retry

import (
	"fmt"
	"github.com/garugaru/resilience4go/resilience"
	"time"
)

type Retry[T any] struct {
	waitStrategy WaitStrategy
	maxRetries   int
}

func New[T any](waitStrategy WaitStrategy, maxRetries int) *Retry[T] {
	return &Retry[T]{
		waitStrategy: waitStrategy,
		maxRetries:   maxRetries,
	}
}

func (r *Retry[T]) Execute(fn resilience.RFunc[T]) resilience.RFunc[T] {
	return func() (T, error) {
		var lastErr error
		for c := 0; c < r.maxRetries; c++ {
			val, err := fn()
			if err != nil {
				lastErr = err
				time.Sleep(r.waitStrategy.Delay())
				r.waitStrategy.OnError()
				continue
			}

			r.waitStrategy.Reset()
			return val, nil
		}
		var empty T
		return empty, fmt.Errorf("retries exhausted after %d times: %w", r.maxRetries, lastErr)
	}
}
