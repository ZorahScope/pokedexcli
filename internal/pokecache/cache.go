package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	stored map[string]cacheEntry
	mu     sync.Mutex
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.stored[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	record, ok := c.stored[key]
	if !ok {
		return []byte{}, false
	}
	return record.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	// defer ticket.Stop()   <-- was here orignally
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			c.mu.Lock()
			for key, entry := range c.stored {
				expirationTime := entry.createdAt.Add(interval)

				if time.Now().After(expirationTime) {
					delete(c.stored, key)
				}
			}
			c.mu.Unlock()
		}
	}()
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	newCache := Cache{
		stored: make(map[string]cacheEntry),
	}
	newCache.reapLoop(interval)
	return &newCache
}
