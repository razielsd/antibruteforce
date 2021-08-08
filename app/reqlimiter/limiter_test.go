package reqlimiter

import (
	"testing"
	"time"

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

func TestReqLimiter_Clean(t *testing.T) {
	reqLimiter := NewReqLimiter(2)
	reqLimiter.items["f1"] = &limiterInfo{
		LastAccess: time.Now().Unix() - limiterTTL - 10,
	}
	reqLimiter.items["f2"] = &limiterInfo{
		LastAccess: time.Now().Unix() - limiterTTL + 10,
	}

	reqLimiter.clean()
	require.Len(t, reqLimiter.items, 1)
	_, ok := reqLimiter.items["f2"]
	require.True(t, ok)
}

func TestReqLimiter_Remove(t *testing.T) {
	key := "test"
	reqLimiter := NewReqLimiter(2)
	require.True(t, reqLimiter.Allow(key))
	require.True(t, reqLimiter.Allow(key))
	require.False(t, reqLimiter.Allow(key))
	reqLimiter.Remove(key)
	require.True(t, reqLimiter.Allow(key))
	require.True(t, reqLimiter.Allow(key))
	require.False(t, reqLimiter.Allow(key))
}
