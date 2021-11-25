package counter

import (
    _ "fmt"
    "strconv"
    "sync"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

// TestHelloWorld
func TestRateLimitReset(t *testing.T) {
    rateLimit := NewRateLimit(5 * time.Second)
    rateLimit.Increment()
    rateLimit.Increment()
    rateLimit.Increment()

    time.Sleep(5 * time.Second)
    rateLimit.Increment()

    assert.Equal(t, uint64(1), rateLimit.counter.Get(), "The counter should be rested and has the count 1.")
}

func TestRateLimitResetParallel(t *testing.T) {
    rateLimit := NewRateLimit(3 * time.Second)

    for c := 0; c < 1000; c++ {
        rateLimit.Increment()
    }
    time.Sleep(4 * time.Second)
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            // time.Sleep(1 * time.Second)
            for c := 0; c < 1000; c++ {
                rateLimit.Increment()
                // fmt.Printf("Count %d\n", count)
            }
            wg.Done()
        }()
    }

    wg.Wait()
    assert.Equal(t, true, rateLimit.reseted, "The counter should be reseted.")
    assert.EqualValues(t, "10000", strconv.FormatUint(rateLimit.counter.Get(), 10), "The counter should be rested and has the count 10000.")
}

func TestRateLimitNextReset(t *testing.T) {
    rateLimit := NewRateLimit(6 * time.Second)
    rateLimit.Increment()
    rateLimit.Increment()
    rateLimit.Increment()

    time.Sleep(3 * time.Second)
    rate := rateLimit.Increment()

    assert.Equal(t, false, rateLimit.reseted, "The counter should be not reseted.")
    assert.Equal(t, uint64(4), rate.Count, "The counter should be 4.")
    assert.Equal(t, int64(3), rate.NextReset, "The counter should be 4.")

    time.Sleep(time.Duration(rate.NextReset) * time.Second)
    rateLimit.Increment()
    assert.Equal(t, true, rateLimit.reseted, "The counter should be reseted.")
}
