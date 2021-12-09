package counter

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
	_ "time"

	"github.com/stretchr/testify/assert"
)

// TestHelloWorld
func TestBucket(t *testing.T) {
	bucket := NewBucket(1 * time.Second)

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func() {
			for c := 0; c < 1000; c++ {
				key := fmt.Sprintf("bucket %d", c%10)
				bucket.Increment(key)
				bucket.Size()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	assert.EqualValues(t, 10, len(bucket.counters), "The bucket should contain 10 counters.")

	assert.Equal(t, false, bucket.Get("bucket 0").reseted, "The counter should not be reseted.")
	assert.EqualValues(t, "1000", strconv.FormatUint(bucket.Get("bucket 0").counter.Get(), 10), "Each counter should be 1000.")
	assert.EqualValues(t, "1000", strconv.FormatUint(bucket.Get("bucket 1").counter.Get(), 10), "Each counter should be 1000.")
	assert.EqualValues(t, "1000", strconv.FormatUint(bucket.Get("bucket 2").counter.Get(), 10), "Each counter should be 1000.")
	assert.EqualValues(t, "1000", strconv.FormatUint(bucket.Get("bucket 3").counter.Get(), 10), "Each counter should be 1000.")
	assert.EqualValues(t, "1000", strconv.FormatUint(bucket.Get("bucket 4").counter.Get(), 10), "Each counter should be 1000.")
	assert.EqualValues(t, "1000", strconv.FormatUint(bucket.Get("bucket 5").counter.Get(), 10), "Each counter should be 1000.")
	assert.EqualValues(t, "1000", strconv.FormatUint(bucket.Get("bucket 6").counter.Get(), 10), "Each counter should be 1000.")
	assert.EqualValues(t, "1000", strconv.FormatUint(bucket.Get("bucket 7").counter.Get(), 10), "Each counter should be 1000.")
	assert.EqualValues(t, "1000", strconv.FormatUint(bucket.Get("bucket 8").counter.Get(), 10), "Each counter should be 1000.")
	assert.EqualValues(t, "1000", strconv.FormatUint(bucket.Get("bucket 9").counter.Get(), 10), "Each counter should be 1000.")
}

func TestBucketWithReset(t *testing.T) {
	bucket := NewBucket(1 * time.Second)

	var wg sync.WaitGroup

	for c := 0; c < 1000; c++ {
		key := fmt.Sprintf("bucket %d", c%10)
		bucket.Increment(key)
	}

	time.Sleep(1 * time.Second)
	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func() {
			for c := 0; c < 1000; c++ {
				key := fmt.Sprintf("bucket %d", c%10)
				bucket.Increment(key)
				bucket.Size()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	bucket.Print()

	assert.Equal(t, int32(10), bucket.Size(), "The bucket should contain 10 counters.")
	assert.Equal(t, "1000", strconv.FormatUint(bucket.Get("bucket 0").counter.Get(), 10), "Each counter should be 1000.")
	assert.Equal(t, "1000", strconv.FormatUint(bucket.Get("bucket 1").counter.Get(), 10), "Each counter should be 1000.")
	assert.Equal(t, "1000", strconv.FormatUint(bucket.Get("bucket 2").counter.Get(), 10), "Each counter should be 1000.")
	assert.Equal(t, "1000", strconv.FormatUint(bucket.Get("bucket 3").counter.Get(), 10), "Each counter should be 1000.")
	assert.Equal(t, "1000", strconv.FormatUint(bucket.Get("bucket 4").counter.Get(), 10), "Each counter should be 1000.")
	assert.Equal(t, "1000", strconv.FormatUint(bucket.Get("bucket 5").counter.Get(), 10), "Each counter should be 1000.")
	assert.Equal(t, "1000", strconv.FormatUint(bucket.Get("bucket 6").counter.Get(), 10), "Each counter should be 1000.")
	assert.Equal(t, "1000", strconv.FormatUint(bucket.Get("bucket 7").counter.Get(), 10), "Each counter should be 1000.")
	assert.Equal(t, "1000", strconv.FormatUint(bucket.Get("bucket 8").counter.Get(), 10), "Each counter should be 1000.")
	assert.Equal(t, "1000", strconv.FormatUint(bucket.Get("bucket 9").counter.Get(), 10), "Each counter should be 1000.")
}

func TestBucketWithRandomCleanup(t *testing.T) {
	bucket := NewBucketWitnCleanup(1 * time.Second)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for c := 0; c < 1000; c++ {
				key := fmt.Sprintf("bucket %d", c%10)
				bucket.Increment(key)
				bucket.Size()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	time.Sleep(2 * time.Second)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for c := 0; c < 1000; c++ {
				key := fmt.Sprintf("bucket %d", c%10)
				bucket.Increment(key)
				bucket.Size()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	time.Sleep(1 * time.Second)

	del := bucket.randomCleanUp(10)
	assert.Equal(t, "10", fmt.Sprintf("%d", del), "All 10 counters should be deleted.")
	assert.Equal(t, 0, len(bucket.counters), "All 10 counters should be deleted.")
	assert.Equal(t, 0, len(bucket.keys), "All 10 counters should be deleted.")

	bucket.Print()
}

func TestBucketSize(t *testing.T) {
	bucket := NewBucket(1 * time.Second)
	
	assert.EqualValues(t, 0, bucket.Size(), "Init bucket size should be an int32 with zero value.")
}
