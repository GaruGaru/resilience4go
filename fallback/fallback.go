package fallback

import "github.com/garugaru/resilience4go/resilience"

type Fallback[T any] struct {
	fallback T
}

func New[T any](fallbackValue T) *Fallback[T] {
	return &Fallback[T]{fallback: fallbackValue}
}

func (f *Fallback[T]) Execute(fn resilience.RFunc[T]) resilience.RFunc[T] {
	return func() (T, error) {
		val, err := fn()
		if err != nil {
			return f.fallback, nil
		}
		return val, nil
	}
}
