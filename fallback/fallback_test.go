package fallback

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFallback_NoErrors(t *testing.T) {
	fb := New("fallback")
	out, err := fb.Execute(func() (string, error) {
		return "ok", nil
	})()
	require.NoError(t, err)
	require.Equal(t, "ok", out)
}

func TestFallback_WithErrors(t *testing.T) {
	fb := New("fallback")
	out, err := fb.Execute(func() (string, error) {
		return "", errors.New("error")
	})()
	require.NoError(t, err)
	require.Equal(t, fb.fallback, out)
}
