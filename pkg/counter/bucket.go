package counter

import (
    "fmt"
    "strconv"
    "sync"
    "time"
)

type Bucket struct {
    counters map[string]*RateLimit
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
    return counter.Increment()
}

func (b *Bucket) Get(key string) *RateLimit {
    return b.counters[key]
}

func (b *Bucket) Size() int {
    return len(b.counters)
}

func (b *Bucket) Print() {
    for k, v := range b.counters {
        fmt.Println(k, "value is", strconv.FormatUint(*v.counter.count, 10))
    }
}
