package ratelimiter

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu      sync.RWMutex
	ttl     uint
	rate    uint
	items   map[string]*BucketHeader
	GetTime func() int64
}

func NewRateLimiter(ttl, rate uint) *RateLimiter {
	return &RateLimiter{
		ttl:   ttl,
		rate:  rate,
		items: make(map[string]*BucketHeader, 0),
		mu:    sync.RWMutex{},
		GetTime: func() int64 {
			return time.Now().Unix()
		},
	}
}

func (r *RateLimiter) Tick(key string) (bool, error) {
	now := r.GetTime()
	hdr := r.getOrCreateHeader(key, now)
	if (hdr.Bucket != nil) && (hdr.Bucket.Time == now) {
		hdr.Bucket.IncCounter()
	} else {
		hdr.appendBucket(now, 1)
	}
	if r.rate < hdr.Count(now, r.ttl) {
		return false, nil
	}
	return true, nil
}

func (r *RateLimiter) getOrCreateHeader(key string, now int64) *BucketHeader {
	r.mu.Lock()
	defer r.mu.Unlock()
	hdr := r.items[key]
	if hdr == nil {
		hdr = NewBucketHeader()
		r.items[key] = hdr
		bucket := NewBucket(now, 0)
		hdr.Bucket = bucket
		hdr.Last = bucket
	}
	return hdr
}
