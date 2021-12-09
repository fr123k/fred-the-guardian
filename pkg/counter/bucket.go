package counter

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fr123k/fred-the-guardian/pkg/utility"
)

type EvictionCallback = func(*Bucket, uint)

func EvictionGarbageCollection(b *Bucket, del uint) {
	threashold := b.randomCleanUpCountSet.SeventyFivePercent()

	if del >= threashold {
		log.Printf("Garbage Collection\n")
		runtime.GC()
	}
}

type Bucket struct {
	counters              map[string]*RateLimit
	keys                  []*string
	size                  atomic.Value
	duration              time.Duration
	mutex                 sync.RWMutex
	randomCleanUpCountGet uint
	randomCleanUpCountSet utility.UInt
	EvictedFunc           EvictionCallback
}

func NewBucket(duration time.Duration) *Bucket {
	bucket := Bucket{
		counters: make(map[string]*RateLimit),
		keys:     make([]*string, 0),
		duration: duration,
	}
	bucket.size.Store(int32(0))
	return &bucket
}

func NewBucketWitnCleanup(duration time.Duration) *Bucket {
	bucket := NewBucket(duration)
	bucket.randomCleanUpCountGet = 1
	bucket.randomCleanUpCountSet = 4
	bucket.EvictedFunc = EvictionGarbageCollection
	return bucket
}

func (b *Bucket) Increment(key string) Rate {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	counter := b.counters[key]
	if counter == nil {
		// TODO perform more random cleanup depending on b.counters map size
		// 4 looks like a good number so that the map size doesn't grow over 170 for
		// max 2 request per second means maximum of 120 counters before expiration
		del := b.randomCleanUp(uint(b.randomCleanUpCountSet))
		b.evicted(del)

		b.counters[key] = NewRateLimit(b.duration)
		counter = b.counters[key]
		b.keys = append(b.keys, &key)
		b.size.Store(int32(len(b.counters)))
		return counter.Increment()
	}
	b.randomCleanUp(b.randomCleanUpCountGet)
	b.size.Store(int32(len(b.counters)))
	return counter.Increment()
}

func (b *Bucket) evicted(del uint) {
	if del > 0 && b.EvictedFunc != nil {
		b.EvictedFunc(b, del)
	}
}

func (b *Bucket) randomCleanUp(maxCnt uint) uint {
	del := uint(0)
	if len(b.keys) <= 0 {
		return del
	}
	rand.Seed(time.Now().UTC().UnixNano())
	for i := uint(0); i < maxCnt; i++ {
		randomIndex := rand.Intn(len(b.keys))
		pick := b.keys[randomIndex]
		v, exists := b.counters[*pick]
		if exists {
			elapsed := time.Now().Sub(v.counter.resetDate)

			if elapsed > v.duration {
				delete(b.counters, *pick)
				//replace the deleted item array position with last item from the array
				b.keys[randomIndex] = b.keys[len(b.keys)-1]
				//shrink the keys array by slicing it with excluding the last item
				b.keys = b.keys[:len(b.keys)-1]
				log.Printf("Deleted Key %d %v, %d %s", randomIndex, *pick, len(b.keys), elapsed)
				del++
				if len(b.keys) <= 0 {
					return del
				}
			}
		}
	}
	return del
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
