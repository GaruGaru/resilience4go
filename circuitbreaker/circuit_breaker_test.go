package circuitbreaker

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_CircuitBreaker(t *testing.T) {
	cb := New[string](
		NewCountBasedWindow(1),
		1, // goes in half open state if the success rate drops below 100%
		1, // pass 100% of requests when in half open state to avoid test undefined behaviour
		0, // immediate state change without time delay ( just request based )
	)

	// execute one function which returns an error, so the circuit goes into CircuitStateOpen state
	_, err := cb.Execute(func() (string, error) {
		return "", errors.New("error")
	})()

	require.Error(t, err)
	require.Equal(t, cb.state, CircuitStateOpen)

	// perform another call to trigger the state change to CircuitStateHalfOpen
	_, err = cb.Execute(func() (string, error) {
		return "", nil
	})()

	// since we are in an CircuitStateOpen state the fail-fast path will return the last error
	require.Error(t, err)
	require.Equal(t, cb.state, CircuitStateHalfOpen)

	// trigger the half open -> closed transition
	_, err = cb.Execute(func() (string, error) {
		return "", nil
	})()
	require.NoError(t, err)
	require.Equal(t, cb.state, CircuitStateClosed)

	// perform the last call to check circuit behaviour in closed state
	val := "success"
	actual, err := cb.Execute(func() (string, error) {
		return val, nil
	})()
	require.NoError(t, err)
	require.Equal(t, val, actual)
	require.Equal(t, cb.state, CircuitStateClosed)

}
