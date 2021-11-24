package counter

import (
    "sync"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

// TestHelloWorld
func TestCounter(t *testing.T) {
    counter := NewCounter()

    assert.EqualValues(t, counter.InitDate().UnixMilli(), counter.ResetDate().UnixMilli(), "Fresh initialized counters init and reset date are equal.")

    var wg sync.WaitGroup

    for i := 0; i < 500; i++ {
        wg.Add(1)
        go func() {
            time.Sleep(1 * time.Second)
            for c := 0; c < 1000; c++ {

                counter.IncrAndGet()
            }
            wg.Done()
        }()
    }

    wg.Wait()

    assert.Equal(t, uint64(500000), counter.Get(), "The two words should be the same.")
}
