package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	cache    map[string]cacheEntry
	mut      *sync.Mutex
	interval time.Duration
	stopChan chan struct{}
}

func NewCache(interval time.Duration) Cache {
	newCache := Cache{
		cache:    make(map[string]cacheEntry),
		mut:      &sync.Mutex{},
		interval: interval,
		stopChan: make(chan struct{}),
	}
	go newCache.reapLoop()
	return newCache
}

func (c *Cache) Add(key string, val []byte) {
	c.mut.Lock()
	defer c.mut.Unlock()
	c.cache[key] = cacheEntry{
		val:       val,
		createdAt: time.Now(),
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mut.Lock()
	defer c.mut.Unlock()
	entry, ok := c.cache[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mut.Lock()
			for key, entry := range c.cache {
				if time.Since(entry.createdAt) > c.interval {
					delete(c.cache, key)
				}
			}
			c.mut.Unlock()
		case <-c.stopChan:
			return
		}
	}
}

func (c *Cache) StopReaping() {
	close(c.stopChan)
}
