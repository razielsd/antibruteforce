package ratelimiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewBucket(t *testing.T) {
	now := time.Now().Unix()
	c := uint(51)
	b := NewBucket(now, c)
	require.NotNil(t, b)
	require.Equal(t, c, b.counter)
	require.Equal(t, now, b.Time)
}

func TestBucket_GetCounter(t *testing.T) {
	now := time.Now().Unix()
	c := uint(90)
	b := NewBucket(now, c)
	require.NotNil(t, b)
	require.Equal(t, c, b.counter)
	require.Equal(t, c, b.GetCounter())
}

func TestBucket_IncCounter(t *testing.T) {
	now := time.Now().Unix()
	c := uint(42)
	b := NewBucket(now, c)
	require.NotNil(t, b)
	b.IncCounter()
	require.Equal(t, c+1, b.GetCounter())
}
