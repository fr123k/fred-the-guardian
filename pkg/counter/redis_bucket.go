package counter

import (
	"context"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisBucket struct {
	name     string
	redis    *redis.Client
	size     atomic.Value
	duration time.Duration
	mutex    sync.RWMutex
}

var ctx = context.Background()

func NewRedisBucket(addr string, duration time.Duration) *RedisBucket {
	bucket := RedisBucket{
		redis: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
		duration: duration,
	}
	bucket.size.Store(int32(0))
	return &bucket
}

func (b *RedisBucket) redisKey(key *string) string {
	return fmt.Sprintf("rate_limit:%s:%s", b.name, *key)
}

func (b *RedisBucket) Increment(key string) Rate {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	rkey := b.redisKey(&key)
	incr := b.redis.Incr(ctx, rkey)
	count, err := incr.Uint64()

	if err != nil {
		panic(err)
	}

	if count <= 1 {
		b.redis.Expire(ctx, rkey, b.duration)
	}

	size := b.redis.DBSize(ctx)
	b.size.Store(int32(size.Val()))
	expire := b.redis.TTL(ctx, rkey)
	return Rate{
		Count:     count,
		NextReset: int64(math.Round(expire.Val().Seconds())),
	}
	// counter := b.counters[key]
	// if counter == nil {
	// 	// TODO perform more random cleanup depending on b.counters map size
	// 	// 4 looks like a good number so that the map size doesn't grow over 170 for
	// 	// max 2 request per second means maximum of 120 counters before expiration
	// 	del := b.randomCleanUp(uint(b.randomCleanUpCountSet))
	// 	b.evicted(del)

	// 	b.counters[key] = NewRateLimit(b.duration)
	// 	counter = b.counters[key]
	// 	b.keys = append(b.keys, &key)
	// 	b.size.Store(int32(len(b.counters)))
	// 	return counter.Increment()
	// }
	// b.randomCleanUp(b.randomCleanUpCountGet)
	// b.size.Store(int32(len(b.counters)))
	// return counter.Increment()
}

// This function is not thread safe.
func (b *RedisBucket) Get(key string) string {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	val := b.redis.Get(ctx, b.redisKey(&key))
	return val.Val()
}

func (b *RedisBucket) Size() int32 {
	return b.size.Load().(int32)
}

// This function is not thread safe.
func (b *RedisBucket) Print() {
	var cursor uint64
	var n int
	for {
		var keys []string
		var err error
		keys, cursor, err = b.redis.Scan(ctx, cursor, "*", 10).Result()
		if err != nil {
			panic(err)
		}
		n += len(keys)
		for _, k := range keys {
			val := b.redis.Get(ctx, k)
			fmt.Println(k, "value is", val.Val())
		}
		if cursor == 0 {
			break
		}
	}

}
