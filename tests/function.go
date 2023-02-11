package tests

import (
	"github.com/garugaru/resilience4go/resilience"
	"sync"
	"time"
)

type countFailingFunction[T any] struct {
	out             T
	err             error
	failuresCount   int
	currentFailures int
	l               *sync.Mutex
}

func CountFailingFn[T any](out T, err error, failuresCount int) resilience.RFunc[T] {
	f := &countFailingFunction[T]{
		out:           out,
		err:           err,
		failuresCount: failuresCount,
		l:             &sync.Mutex{},
	}
	return f.Function()
}

func (c *countFailingFunction[T]) Function() resilience.RFunc[T] {
	return func() (T, error) {
		c.l.Lock()
		defer c.l.Unlock()
		if c.currentFailures <= c.failuresCount {
			c.currentFailures++
			return c.out, c.err
		}
		return c.out, nil
	}
}

type timeFailingFunction[T any] struct {
	out       T
	err       error
	failAfter time.Duration
	firstCall *time.Time
	l         *sync.Mutex
}

func TimeFailingFn[T any](out T, err error, failAfter time.Duration) resilience.RFunc[T] {
	f := &timeFailingFunction[T]{
		out:       out,
		err:       err,
		failAfter: failAfter,
		l:         &sync.Mutex{},
	}
	return f.Function()
}

func (c *timeFailingFunction[T]) Function() resilience.RFunc[T] {
	return func() (T, error) {
		c.l.Lock()
		defer c.l.Unlock()
		if c.firstCall == nil {
			t := time.Now()
			c.firstCall = &t
			return c.out, c.err
		}

		if time.Since(*c.firstCall) < c.failAfter {
			return c.out, c.err
		}

		return c.out, nil
	}
}
