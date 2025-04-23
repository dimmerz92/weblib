package weblib

import (
	"sync"
	"time"
)

type cacheItem struct {
	lastAccess time.Time
	value      any
}

type Cache struct {
	data     map[string]*cacheItem
	mu       sync.Mutex
	ttl      time.Duration
	interval time.Duration
	close    chan bool
	once     sync.Once
}

// Close signals for the cache to gracefully stop the goroutine that periodically cleans expired data.
func (c *Cache) Close() {
	c.once.Do(func() {
		c.close <- true
		close(c.close)
	})
}

// cleaner is a goroutine function that removes expired data from the cache at the specified cache clean interval.
func (c *Cache) cleaner() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-c.close:
			return

		case now := <-ticker.C:
			c.mu.Lock()
			for k, v := range c.data {
				if now.Sub(v.lastAccess) > c.ttl {
					delete(c.data, k)
				}
			}
			c.mu.Unlock()
		}
	}
}

// NewCache returns a new in memory only cache that self cleans expired data at the given cleanInterval.
// The expiry of an item in the cache is determined by the ttl.
func NewCache(ttl, cleanInterval time.Duration) *Cache {
	cache := &Cache{
		data:     make(map[string]*cacheItem),
		mu:       sync.Mutex{},
		ttl:      ttl,
		interval: cleanInterval,
		close:    make(chan bool),
	}

	// start the cleaning goroutine
	go cache.cleaner()

	return cache
}

// Put adds the value v to the cache with key k.
func (c *Cache) Put(k string, v any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[k] = &cacheItem{value: v, lastAccess: time.Now()}
}

// Get retrieves the value v from the cache with key k if it exists and is not expired.
func (c *Cache) Get(k string) (v any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, ok := c.data[k]
	if ok && time.Since(data.lastAccess) <= c.ttl {
		v = data.value
		data.lastAccess = time.Now()
	}

	return
}

// Delete removes the key value pair mapped to by k if it exists.
func (c *Cache) Delete(k string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, k)
}
