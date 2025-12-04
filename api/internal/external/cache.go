package external

import (
	"encoding/json"
	"sync"
	"time"
)

// cache provides a simple in-memory cache with TTL for external API responses.
type cache struct {
	mu      sync.RWMutex
	entries map[string]cacheEntry
}

type cacheEntry struct {
	value     json.RawMessage
	expiresAt time.Time
}

func newCache() *cache {
	return &cache{
		entries: make(map[string]cacheEntry),
	}
}

// get returns a cached value if it exists and has not expired.
func (c *cache) get(key string) (json.RawMessage, bool) {
	c.mu.RLock()
	entry, ok := c.entries[key]
	c.mu.RUnlock()
	if !ok {
		return nil, false
	}
	if time.Now().After(entry.expiresAt) {
		c.mu.Lock()
		delete(c.entries, key)
		c.mu.Unlock()
		return nil, false
	}
	return entry.value, true
}

// set stores a value with the provided TTL.
func (c *cache) set(key string, value json.RawMessage, ttl time.Duration) {
	c.mu.Lock()
	c.entries[key] = cacheEntry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
	c.mu.Unlock()
}
