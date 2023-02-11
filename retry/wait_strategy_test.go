package retry

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestWaitStrategy_Fixed(t *testing.T) {
	delay := 10 * time.Millisecond
	s := NewFixedDelay(delay)
	s.Reset()
	s.OnError()
	require.Equal(t, delay, s.Delay())
}

func TestWaitStrategy_Exponential_No_Errors(t *testing.T) {
	delay := 10 * time.Millisecond
	s := NewExponentialBackoff(delay, 1*time.Second)
	s.Reset()
	require.Equal(t, delay, s.Delay())
}

func TestWaitStrategy_Exponential_WithErrors(t *testing.T) {
	delay := 10 * time.Millisecond
	s := NewExponentialBackoff(delay, 1*time.Second)
	s.OnError()
	s.OnError()
	s.OnError()
	require.Equal(t, 80*time.Millisecond, s.Delay())
}

func TestWaitStrategy_Exponential_WithErrors_Reset(t *testing.T) {
	delay := 10 * time.Millisecond
	s := NewExponentialBackoff(delay, 1*time.Second)
	s.OnError()
	s.Reset()
	require.Equal(t, delay, s.Delay())
}

func TestWaitStrategy_Exponential_Max_Cap(t *testing.T) {
	max := 1 * time.Second
	s := NewExponentialBackoff(10*time.Millisecond, max)
	for i := 0; i < 50; i++ {
		s.OnError()
	}
	require.Equal(t, max, s.Delay())
}
