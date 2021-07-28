package ratelimiter

import (
	"sync"
)

type Bucket struct {
	mu      sync.RWMutex
	Time    int64
	Prev    *Bucket
	Next    *Bucket
	counter uint
}

func NewBucket(now int64, counter uint) *Bucket {
	return &Bucket{
		Time:    now,
		counter: counter,
		mu:      sync.RWMutex{},
	}
}

func (b *Bucket) IncCounter() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.counter++
}

func (b *Bucket) GetCounter() uint {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.counter
}
