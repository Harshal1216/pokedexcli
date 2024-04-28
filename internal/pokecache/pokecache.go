package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	cacheMap map[string]cacheEntry
	lock     *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(cacheInterval time.Duration) Cache {
	cache := Cache{
		cacheMap: make(map[string]cacheEntry),
		lock:     &sync.Mutex{},
	}
	ticker := time.NewTicker(cacheInterval)
	go cache.reapLoop(ticker, cacheInterval)
	return cache
}

func (c *Cache) Add(key string, value []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cacheMap[key] = cacheEntry{
		createdAt: time.Now(),
		val:       value,
	}
}

func (c *Cache) Get(key string) ([]byte, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	value, ok := c.cacheMap[key]
	if !ok {
		return nil, fmt.Errorf("cache miss for key: %s", key)
	}
	return value.val, nil
}

func (c *Cache) reapLoop(ticker *time.Ticker, cacheInterval time.Duration) {
	// Listen to tick on channel
	for range ticker.C {
		c.lock.Lock()
		currentTime := time.Now()
		for key, val := range c.cacheMap {
			if currentTime.Sub(val.createdAt) > cacheInterval {
				delete(c.cacheMap, key)
			}
		}
		c.lock.Unlock()
	}
}
