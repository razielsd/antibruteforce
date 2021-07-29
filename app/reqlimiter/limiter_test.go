package reqlimiter

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewRateLimiter(t *testing.T) {
	var ttl time.Duration = 10
	rate := 200
	l := NewReqLimiter(ttl, rate)
	require.NotNil(t, l)
	require.Equal(t, ttl, l.ttl)
	require.Equal(t, rate, l.rate)
}

func TestReqLimiter_Allow(t *testing.T) {
	reqLimiter := NewReqLimiter(60, 2)
	require.True(t, reqLimiter.Allow("test"))
	require.True(t, reqLimiter.Allow("test"))
	require.False(t, reqLimiter.Allow("test"))
	require.True(t, reqLimiter.Allow("test2"))
}
