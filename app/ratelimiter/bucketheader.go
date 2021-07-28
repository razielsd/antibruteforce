package ratelimiter

type BucketHeader struct {
	Bucket *Bucket
	Last   *Bucket
}

func NewBucketHeader() *BucketHeader {
	return &BucketHeader{}
}

func (h *BucketHeader) Count(now int64, ttl uint) uint {
	b := h.Bucket
	var c uint
	expired := now - int64(ttl)
	for b != nil {
		if b.Time > expired {
			c += b.GetCounter()
			b = b.Prev
		} else {
			b = nil
		}
	}
	return c
}

func (h *BucketHeader) appendBucket(now int64, counter uint) {
	bucket := NewBucket(now, counter)
	bucket.Prev = h.Bucket
	if h.Bucket != nil {
		h.Bucket.Next = bucket
	}
	h.Bucket = bucket
}
