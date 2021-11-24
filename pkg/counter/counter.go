package counter

import (
    "sync/atomic"
    "time"
)

type Counter struct {
    resetDate time.Time
    initDate  time.Time
    count     *uint64
}

func NewCounter() *Counter {
    now := time.Now()
    zero := uint64(0)
    return &Counter{
        resetDate: now,
        initDate:  now,
        count:     &zero,
    }
}

func (c *Counter) IncrAndGet() uint64 {
    return atomic.AddUint64(c.count, 1)
}

func (c *Counter) Get() uint64 {
    return atomic.LoadUint64(c.count)
}

func (c *Counter) InitDate() time.Time {
    return c.initDate
}

func (c *Counter) ResetDate() time.Time {
    return c.resetDate
}

func (c *Counter) Reset() {
    c.resetDate = time.Now()
    atomic.StoreUint64(c.count, uint64(0))
}
