package counter

import (
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type Bucket struct {
	counters map[string]*RateLimit
	size     atomic.Value
	duration time.Duration
	mutex    sync.RWMutex
}

func NewBucket(duration time.Duration) *Bucket {
	return &Bucket{
		counters: make(map[string]*RateLimit),
		duration: duration,
	}
}

func (b *Bucket) Increment(key string) Rate {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	counter := b.counters[key]
	if counter == nil {
		b.counters[key] = NewRateLimit(b.duration)
		counter = b.counters[key]
	}
	b.size.Store(int32(len(b.counters)))
	return counter.Increment()
}

// This function is not thread safe.
func (b *Bucket) Get(key string) *RateLimit {
	return b.counters[key]
}

func (b *Bucket) Size() int32 {
	return b.size.Load().(int32)
}

// This function is not thread safe.
func (b *Bucket) Print() {
	for k, v := range b.counters {
		fmt.Println(k, "value is", strconv.FormatUint(*v.counter.count, 10))
	}
}
