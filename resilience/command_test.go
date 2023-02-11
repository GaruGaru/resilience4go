package resilience

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

type MockMiddleware[T any] struct {
	fn func(fn RFunc[T]) RFunc[T]
}

func (m *MockMiddleware[T]) Execute(r RFunc[T]) RFunc[T] {
	return m.fn(r)
}

func TestCommand_Run(t *testing.T) {
	middlewareStack := make([]int, 0)
	m0 := &MockMiddleware[string]{
		fn: func(fn RFunc[string]) RFunc[string] {
			return func() (string, error) {
				middlewareStack = append(middlewareStack, 0)
				return fn()
			}
		},
	}
	m1 := &MockMiddleware[string]{
		fn: func(fn RFunc[string]) RFunc[string] {
			return func() (string, error) {
				middlewareStack = append(middlewareStack, 1)
				return fn()
			}
		},
	}

	var runner = New[string](m0, m1)

	_, err := runner.Execute(func() (string, error) {
		return "", errors.New("error")
	})

	require.Error(t, err)
}
