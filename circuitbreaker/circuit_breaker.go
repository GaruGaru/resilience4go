package circuitbreaker

import (
	"fmt"
	"github.com/garugaru/resilience4go/resilience"
	"math/rand"
	"time"
)

type CircuitState int

const (
	CircuitStateOpen CircuitState = iota
	CircuitStateClosed
	CircuitStateHalfOpen
)

type CircuitBreaker[T any] struct {
	window                Window
	state                 CircuitState
	successRateThreshold  float64
	halfOpenRequestsRatio float64
	closedDuration        time.Duration
	lastStateTransition   time.Time
	stateStay             time.Duration
	lastErr               error
}

func New[T any](window Window, successRateThreshold float64, halfOpenRequestsRatio float64, stateStay time.Duration) *CircuitBreaker[T] {
	return &CircuitBreaker[T]{
		window:                window,
		state:                 CircuitStateClosed,
		successRateThreshold:  successRateThreshold,
		closedDuration:        stateStay,
		stateStay:             stateStay,
		halfOpenRequestsRatio: halfOpenRequestsRatio,
	}
}

func (c *CircuitBreaker[T]) Execute(fn resilience.RFunc[T]) resilience.RFunc[T] {
	return func() (T, error) {
		var empty T
		if c.state == CircuitStateOpen {
			if time.Since(c.lastStateTransition) > c.closedDuration {
				c.transitionToState(CircuitStateHalfOpen)
			}
			return empty, fmt.Errorf("circuit open: %w", c.lastErr)
		}

		if c.state == CircuitStateHalfOpen {
			if rand.Float64() > c.halfOpenRequestsRatio {
				return empty, fmt.Errorf("circuit open: %w", c.lastErr)
			}
		}

		res, err := fn()
		if err != nil {
			c.lastErr = err
			c.window.Notify(WindowCallFailed)

			if c.canChangeState() {
				if c.state == CircuitStateHalfOpen || c.state == CircuitStateClosed {
					if c.window.SuccessRate() < c.successRateThreshold {
						c.transitionToState(CircuitStateOpen)
					}
				}
			}

		} else {
			c.window.Notify(WindowCallSuccess)

			if c.canChangeState() {
				if c.state == CircuitStateHalfOpen {
					if c.window.SuccessRate() >= c.successRateThreshold {
						c.transitionToState(CircuitStateClosed)
					}
				}
			}
		}

		return res, err
	}
}

func (c *CircuitBreaker[T]) transitionToState(newState CircuitState) {
	c.state = newState
	c.lastStateTransition = time.Now()
}

func (c *CircuitBreaker[T]) canChangeState() bool {
	return time.Since(c.lastStateTransition) > c.stateStay
}
