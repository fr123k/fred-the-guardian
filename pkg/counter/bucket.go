package counter

import (
	"log"
	"runtime"
)

type EvictionCallback = func(*MemoryBucket, uint)

func EvictionGarbageCollection(b *MemoryBucket, del uint) {
	threashold := b.randomCleanUpCountSet.SeventyFivePercent()

	if del >= threashold {
		log.Printf("Garbage Collection\n")
		runtime.GC()
	}
}

type Bucket interface {
	Size() int32
	Increment(string) Rate
}
