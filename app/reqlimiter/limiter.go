package reqlimiter

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const limiterTTL = 60

type ReqLimiter struct {
	mu    sync.RWMutex
	ttl   time.Duration
	rate  int
	items map[string]*rate.Limiter
}

func NewReqLimiter(reqRate int) *ReqLimiter {
	return &ReqLimiter{
		ttl:   limiterTTL,
		rate:  reqRate,
		items: make(map[string]*rate.Limiter),
		mu:    sync.RWMutex{},
	}
}

func (r *ReqLimiter) Allow(key string) bool {
	limiter := r.getOrCreateLimiter(key)
	return limiter.Allow()
}

func (r *ReqLimiter) getOrCreateLimiter(key string) *rate.Limiter {
	r.mu.Lock()
	defer r.mu.Unlock()
	limiter := r.items[key]
	if limiter == nil {
		limiter = rate.NewLimiter(rate.Every(r.ttl*time.Second), r.rate)
		r.items[key] = limiter
	}
	return limiter
}
