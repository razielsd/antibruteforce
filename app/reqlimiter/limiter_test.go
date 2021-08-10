package reqlimiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const testKey = "test"

func TestNewRateLimiter(t *testing.T) {
	rate := 200
	l := NewReqLimiter(NewLimiterConfig(rate))
	require.NotNil(t, l)
	require.Equal(t, rate, l.rate)
}

func TestReqLimiter_Allow(t *testing.T) {
	reqLimiter := NewReqLimiter(NewLimiterConfig(2))
	require.True(t, reqLimiter.Allow("test"))
	require.True(t, reqLimiter.Allow("test"))
	require.False(t, reqLimiter.Allow("test"))
	require.True(t, reqLimiter.Allow("test2"))
}

func TestReqLimiter_Clean(t *testing.T) {
	reqLimiter := NewReqLimiter(NewLimiterConfig(2))
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
	reqLimiter := NewReqLimiter(NewLimiterConfig(2))
	require.True(t, reqLimiter.Allow(testKey))
	require.True(t, reqLimiter.Allow(testKey))
	require.False(t, reqLimiter.Allow(testKey))
	reqLimiter.Remove(testKey)
	require.True(t, reqLimiter.Allow(testKey))
	require.True(t, reqLimiter.Allow(testKey))
	require.False(t, reqLimiter.Allow(testKey))
}

func TestReqLimiter_CleanByTimer(t *testing.T) {
	cfg := NewLimiterConfig(2)
	cfg.TTL = 1
	cfg.CleanInterval = 1010 * time.Millisecond
	reqLimiter := NewReqLimiter(cfg)
	require.True(t, reqLimiter.Allow(testKey))
	require.True(t, reqLimiter.HasKey(testKey))
	cond := func() bool {
		return reqLimiter.HasKey(testKey)
	}
	require.Eventually(t, cond, 5*time.Second, 100*time.Millisecond)
}

func TestReqLimiter_TokenExpire(t *testing.T) {
	cfg := NewLimiterConfig(2)
	cfg.TTL = 1
	reqLimiter := NewReqLimiter(cfg)
	require.True(t, reqLimiter.Allow(testKey))
	require.True(t, reqLimiter.Allow(testKey))
	time.Sleep(cfg.TTL * time.Second)
	require.True(t, reqLimiter.Allow(testKey))
}

func TestReqLimiter_HasKey_Exists(t *testing.T) {
	reqLimiter := NewReqLimiter(NewLimiterConfig(2))
	require.True(t, reqLimiter.Allow(testKey))
	require.True(t, reqLimiter.HasKey(testKey))
}

func TestReqLimiter_HasKey_NotExists(t *testing.T) {
	reqLimiter := NewReqLimiter(NewLimiterConfig(2))
	require.False(t, reqLimiter.HasKey(testKey))
}
