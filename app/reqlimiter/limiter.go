package reqlimiter

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	limiterTTL           = 60
	limiterCleanInterval = 61
)

type ReqLimiter struct {
	mu    sync.Mutex
	ttl   time.Duration
	rate  int
	items map[string]*limiterInfo
}

type limiterInfo struct {
	Limiter    *rate.Limiter
	LastAccess int64
	mu         sync.Mutex
}

type LimiterConfig struct {
	TTL           time.Duration
	CleanInterval time.Duration
	Rate          int
}

// NewReqLimiter create new instance of ReqLimiter.
func NewReqLimiter(cfg LimiterConfig) *ReqLimiter {
	limiter := &ReqLimiter{
		ttl:   cfg.TTL,
		rate:  cfg.Rate,
		items: make(map[string]*limiterInfo),
		mu:    sync.Mutex{},
	}
	go func() {
		for range time.Tick(cfg.CleanInterval) {
			limiter.clean()
		}
	}()
	return limiter
}

// NewLimiterConfig create default limiter configuration.
func NewLimiterConfig(rate int) LimiterConfig {
	return LimiterConfig{
		Rate:          rate,
		TTL:           limiterTTL,
		CleanInterval: limiterCleanInterval * time.Second,
	}
}

// Allow check allow or not action by limits.
func (r *ReqLimiter) Allow(key string) bool {
	li := r.getOrCreateLimiter(key)
	return li.Allow()
}

func (r *ReqLimiter) getOrCreateLimiter(key string) *limiterInfo {
	r.mu.Lock()
	defer r.mu.Unlock()
	li, ok := r.items[key]
	if !ok {
		limiter := rate.NewLimiter(rate.Every(r.ttl*time.Second), r.rate)
		li = &limiterInfo{
			Limiter: limiter,
		}
		r.items[key] = li
	}
	return li
}

// Remove remove bucket.
func (r *ReqLimiter) Remove(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.items, key)
}

// Allow check allow or not action by limits.
func (l *limiterInfo) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.LastAccess = time.Now().Unix()
	return l.Limiter.Allow()
}

func (r *ReqLimiter) clean() {
	r.mu.Lock()
	defer r.mu.Unlock()
	expire := time.Now().Unix() - int64(r.ttl)
	for i, v := range r.items {
		if expire > v.LastAccess {
			delete(r.items, i)
		}
	}
}
