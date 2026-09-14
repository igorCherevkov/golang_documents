package cache

import (
	"sync"
	"time"
)

type Record struct {
	Status	int
	ContentType	string
	Body	[]byte
	ExpiresAt	time.Time
}

type Cache struct {
	mu	sync.RWMutex
	records	map[string]Record
}

func NewCache() *Cache {
	return &Cache{
		records: make(map[string]Record),
	}
}

func (c *Cache) Get(key string) (Record, bool) {
	c.mu.RLock()
	record, ok := c.records[key]
	c.mu.RUnlock()

	if !ok {
		return Record{}, false
	}

	if record.ExpiresAt.IsZero() && time.Now().After(record.ExpiresAt) {
		return Record{}, false
	}
 
	return record, true
}

func (c *Cache) Set(key string, record Record, ttl time.Duration) {
	if ttl > 0 {
		record.ExpiresAt = time.Now().Add(ttl)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.records[key] = record
}

func (c * Cache) InvalidatePrefix(prefix string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for record := range c.records {
		if len(record) >= len(prefix) && record[:len(prefix)] == prefix {
			delete(c.records, record)
		}
	}
}