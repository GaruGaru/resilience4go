package retry

import (
	"errors"
	"fmt"
	"github.com/garugaru/resilience4go/tests"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRetry_NoErrors(t *testing.T) {
	rfn := New[string](NewFixedDelay(10*time.Millisecond), 10)
	_, err := rfn.Execute(func() (string, error) {
		return "", nil
	})()

	require.NoError(t, err)
}

func TestRetry_AfterN_Errors(t *testing.T) {
	var expectedOut = "test"
	rfn := New[string](NewFixedDelay(10*time.Millisecond), 10)
	fn := tests.CountFailingFn[string](expectedOut, errors.New("error"), 5)

	out, err := rfn.Execute(fn)()
	require.NoError(t, err)
	require.Equal(t, expectedOut, out)
}

func TestRetry_AfterN_Errors_Retries_Exhausted(t *testing.T) {
	var expectedErr = errors.New("error")
	rfn := New[string](NewFixedDelay(10*time.Millisecond), 3)
	fn := tests.CountFailingFn[string]("", expectedErr, 5)

	_, err := rfn.Execute(fn)()
	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
}

func TestRetry_After_Transient_Errors(t *testing.T) {
	var expectedOut = "test"
	rfn := New[string](NewFixedDelay(100*time.Millisecond), 3)
	fn := tests.TimeFailingFn[string](expectedOut, errors.New("error"), 80*time.Millisecond)

	out, err := rfn.Execute(fn)()
	require.NoError(t, err)
	require.Equal(t, expectedOut, out)
}

func TestRetry_Propagate_Error(t *testing.T) {
	var (
		expectedOut = "test"
		werr        = errors.New("error")
	)
	rfn := New[string](NewFixedDelay(10*time.Millisecond), 1)
	fn := tests.CountFailingFn[string](expectedOut, fmt.Errorf("generic error: %w", werr), 5)

	_, err := rfn.Execute(fn)()
	require.Error(t, err)
	require.ErrorIs(t, err, werr)
}
