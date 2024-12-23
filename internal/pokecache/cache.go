package pokecache

import (
	"sync"
	"time"
)

type cache struct {
	stored map[string]cacheEntry
	mu     sync.Mutex
}

func (c *cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.stored[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	record, ok := c.stored[key]
	if !ok {
		return []byte{}, false
	}
	return record.val, true
}

func (c *cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	deleteExpiredCache := func(ch<- chan time.Time) {
		//receive time from channel
		//loop through cache
		//delete anything older than interval
	}
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *cache {
	newCache := cache{
		stored: make(map[string]cacheEntry),
	}
	newCache.reapLoop(interval)
	return &newCache
}
