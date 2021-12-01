package counter

import (
    "math"
    "sync"
    "time"
)

type RateLimit struct {
    counter  *Counter
    duration time.Duration
    mutex    sync.Mutex
    reseted  bool
}

type Rate struct {
    Count     uint64
    NextReset int64
}

func NewRateLimit(duration time.Duration) *RateLimit {
    return &RateLimit{
        counter:  NewCounter(),
        duration: duration,
    }
}

func (r *RateLimit) Increment() Rate {
    start := time.Now()
    elapsed := start.Sub(r.counter.resetDate)
    nextReset := r.duration - elapsed

    if elapsed > r.duration {
        // TODO add proper logging
        // fmt.Printf("Elapsed time %s > %s \n", start.String(), r.counter.resetDate.String())
        r.mutex.Lock()
        defer r.mutex.Unlock()
        elapsed = start.Sub(r.counter.resetDate)
        // double check idiom to ensure only lock when its nessary and to only reset the counter once
        if elapsed > r.duration {
            // TODO add proper logging
            // fmt.Printf("Reset %s > %s \n", start.String(), r.counter.resetDate.String())
            r.counter.Reset()
            r.reseted = true
        }
        nextReset = r.duration - elapsed
    }
    return Rate{
        Count:     r.counter.IncrAndGet(),
        NextReset: int64(math.Round(nextReset.Seconds())),
    }
}
