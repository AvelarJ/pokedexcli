package internal

import (
	//"fmt"
	"sync"
	"time"
)

type Cache struct {
	content map[string]cacheEntry
	mutex   sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

// init a new cache with a specified interval for reaping expired entries
func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		content: make(map[string]cacheEntry),
	}
	go cache.reapLoop(interval)
	return cache
}

// Cache methods

// reapLoop runs in a separate goroutine and periodically checks for expired entries in the cache,
// removing them if they have been in the cache longer than the specified interval.
func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		c.mutex.Lock()
		for key, entry := range c.content {
			if time.Since(entry.createdAt) > interval {
				delete(c.content, key)
			}
		}
		c.mutex.Unlock()
	}
}

// Add adds a new entry to the cache with the given key and value.
func (c *Cache) Add(key string, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.content[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

// Get retrieves the value associated with the given key if it exists, nil and false otherwise.
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	entry, exits := c.content[key]
	if !exits {
		return nil, false
	}
	return entry.val, true
}
