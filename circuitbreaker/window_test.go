package circuitbreaker

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWindow_Count(t *testing.T) {
	wc := NewCountBasedWindow(3)
	wc.Notify(WindowCallFailed)
	wc.Notify(WindowCallSuccess)
	wc.Notify(WindowCallFailed)

	require.Equal(t, 1.0/3.0, wc.SuccessRate())
}

func TestWindow_Count_All_Failures(t *testing.T) {
	wc := NewCountBasedWindow(3)
	wc.Notify(WindowCallFailed)
	wc.Notify(WindowCallFailed)
	wc.Notify(WindowCallFailed)

	require.Equal(t, 0.0, wc.SuccessRate())
}

func TestWindow_Count_All_Success(t *testing.T) {
	wc := NewCountBasedWindow(3)
	wc.Notify(WindowCallSuccess)
	wc.Notify(WindowCallSuccess)
	wc.Notify(WindowCallSuccess)

	require.Equal(t, 1.0, wc.SuccessRate())
}
