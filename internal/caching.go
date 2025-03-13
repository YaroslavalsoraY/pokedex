package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	sync.RWMutex
	Cached 	 map[string]cacheEntry
	Interval time.Duration
}

type cacheEntry struct {
	CreatedAt time.Time
	Val       []byte
}

func NewCache(interval time.Duration) *Cache{
	items := make(map[string]cacheEntry)

	cache := Cache{
		Cached: 	 items,
		Interval: interval,
	}
	if interval > 0 {
		go cache.reapLoop()
	}
	return &cache
}

func (c *Cache) Add(key string, value []byte) {
	c.Lock()

	defer c.Unlock()

	c.Cached[key] = cacheEntry{
		CreatedAt: time.Now(),
		Val:       value,
	}
}

func (c *Cache) Get(key string) ([]byte, bool){
	c.Lock()

	defer c.Unlock()

	_, ok := c.Cached[key]
	if !ok {
		return nil, false
	}
	return c.Cached[key].Val, true
}

func (c *Cache) reapLoop() {
	for {
		<-time.After(c.Interval * time.Second)
		if c.Cached == nil {
			return
		}

		for k, v := range c.Cached {
			if v.CreatedAt.Sub(time.Now()) * time.Second > c.Interval * time.Second { 
				c.Lock()
				delete(c.Cached, k)
				c.Unlock()
			}
		}	
	}
}