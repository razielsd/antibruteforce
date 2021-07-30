package reqlimiter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRateLimiter(t *testing.T) {
	rate := 200
	l := NewReqLimiter(rate)
	require.NotNil(t, l)
	require.Equal(t, rate, l.rate)
}

func TestReqLimiter_Allow(t *testing.T) {
	reqLimiter := NewReqLimiter(2)
	require.True(t, reqLimiter.Allow("test"))
	require.True(t, reqLimiter.Allow("test"))
	require.False(t, reqLimiter.Allow("test"))
	require.True(t, reqLimiter.Allow("test2"))
}
