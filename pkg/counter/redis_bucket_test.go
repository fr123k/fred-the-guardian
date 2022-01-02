//go:build integration

package counter

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type CleanUp func()

func RedisContainer(t *testing.T) (*dockertest.Resource, CleanUp) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		panic(err)
	}

	resource, err := pool.Run("redis", "latest", []string{})
	if err != nil {
		panic(err)
	}

	var redisInfo string

	err = pool.Retry(func() error {
		client := redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%s", "localhost", resource.GetPort("6379/tcp")),
		})

		redisInfo, err = client.Info(ctx).Result()
		if err != nil {
			t.Log("container not ready, waiting...")
			return err
		}
		return nil
	})

	require.NoError(t, err, "HTTP error")

	require.Contains(t, redisInfo, "redis_mode:standalone", "does not respond with love?")

	return resource, func() {
		require.NoError(t, pool.Purge(resource), "failed to remove container")
	}
}

func TestRespondsWithLove(t *testing.T) {

	_, cleanUpFunc := RedisContainer(t)

	t.Cleanup(cleanUpFunc)
}

func TestRedisCounter(t *testing.T) {
	resource, cleanUpFunc := RedisContainer(t)

	t.Cleanup(cleanUpFunc)

	bucket := NewRedisBucket(fmt.Sprintf("%s:%s", "localhost", resource.GetPort("6379/tcp")), 2 * time.Second)
	
	var wg sync.WaitGroup

	for c := 0; c < 1000; c++ {
		key := fmt.Sprintf("bucket %d", c%10)
		bucket.Increment(key)
	}

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
	bucket.Print()

	assert.EqualValues(t, 10, bucket.Size(), "The bucket should contain 10 counters.")

	assert.EqualValues(t, "1000", bucket.Get("bucket 0"), "counter bucket 0 should be 1000.")
	assert.EqualValues(t, "1000", bucket.Get("bucket 1"), "counter bucket 1 should be 1000.")
	assert.EqualValues(t, "1000", bucket.Get("bucket 2"), "counter bucket 2 should be 1000.")
	assert.EqualValues(t, "1000", bucket.Get("bucket 3"), "counter bucket 3 should be 1000.")
	assert.EqualValues(t, "1000", bucket.Get("bucket 4"), "counter bucket 4 should be 1000.")
	assert.EqualValues(t, "1000", bucket.Get("bucket 5"), "counter bucket 5 should be 1000.")
	assert.EqualValues(t, "1000", bucket.Get("bucket 6"), "counter bucket 6 should be 1000.")
	assert.EqualValues(t, "1000", bucket.Get("bucket 7"), "counter bucket 7 should be 1000.")
	assert.EqualValues(t, "1000", bucket.Get("bucket 8"), "counter bucket 8 should be 1000.")
	assert.EqualValues(t, "1000", bucket.Get("bucket 9"), "counter bucket 9 should be 1000.")

	bucket.Print()
}
