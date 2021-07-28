package ratelimiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewBucketHeader(t *testing.T) {
	h := NewBucketHeader()
	require.NotNil(t, h)
}

func TestBucketHeader_Count(t *testing.T) {
	tests := []struct {
		name   string
		count  int
		ttl    uint
		bucket []int
	}{
		{
			name:   "empty bucket",
			count:  0,
			ttl:    20,
			bucket: []int{0},
		},
		{
			name:   "single bucket",
			count:  10,
			ttl:    20,
			bucket: []int{10},
		},
		{
			name:   "multi bucket",
			count:  35,
			ttl:    20,
			bucket: []int{10, 5, 15, 0, 0, 0, 5},
		},
		{
			name:   "multi bucket with expired",
			count:  30,
			ttl:    5,
			bucket: []int{10, 5, 15, 0, 0, 0, 5},
		},
		{
			name:   "all bucket expired",
			count:  0,
			ttl:    5,
			bucket: []int{0, 0, 0, 0, 0, 0, 5, 10},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := NewBucketHeader()
			now := time.Now().Unix()
			appendHeaderBucket(h, now, test.bucket)
			require.Equal(t, uint(test.count), h.Count(now, test.ttl))
		})
	}
}

func appendHeaderBucket(h *BucketHeader, now int64, data []int) {
	l := len(data)
	for i := range data {
		v := data[l-i-1]
		if v > 0 {
			h.appendBucket(now-int64(l-i-1), uint(v))
		}
	}
}
