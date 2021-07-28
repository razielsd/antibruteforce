package ratelimiter

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewRateLimiter(t *testing.T) {
	ttl := uint(60)
	rate := uint(200)
	l := NewRateLimiter(ttl, rate)
	require.NotNil(t, l)
	require.Equal(t, ttl, l.ttl)
	require.Equal(t, rate, l.rate)
}

func TestRateLimiter_Tick_NoBucketExceed(t *testing.T) {
	rate := 3
	l := NewRateLimiter(uint(10), uint(rate))
	for i := 0; i < rate; i++ {
		l.GetTime = func() int64 { return 100 + int64(i) }
		ok, err := l.Tick("ip")
		require.NoError(t, err)
		require.True(t, ok)
	}

	ok, err := l.Tick("ip")
	require.NoError(t, err)
	require.False(t, ok)

	ok, err = l.Tick("ip2")
	require.NoError(t, err)
	require.True(t, ok)
}

func TestRateLimiter_Tick_BucketExceed(t *testing.T) {
	rate := 10
	l := NewRateLimiter(uint(10), uint(rate))
	i := 0
	for i = 0; i < rate+5; i++ {
		l.GetTime = func() int64 {
			i := i
			return 100 + int64(i)
		}
		ok, err := l.Tick("ip")
		require.NoError(t, err)
		require.True(t, ok, fmt.Sprintf("i=%d", i))
	}
	i--
	ok, err := l.Tick("ip")
	require.NoError(t, err)
	require.False(t, ok)

	ok, err = l.Tick("ip2")
	require.NoError(t, err)
	require.True(t, ok)
}

func TestRateLimiter_GetTime(t *testing.T) {
	l := NewRateLimiter(uint(10), uint(10))
	require.InDelta(t, time.Now().Unix(), l.GetTime(), 1)
}
