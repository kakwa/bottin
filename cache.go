package bottin

import (
	"encoding/json"
	"sync"
	"time"
)

// Cache stores slices of RR structs, with each key mapping to an RR slice.
type Cache struct {
	items map[string][]RR
	mutex sync.RWMutex
}

// NewCache initializes a new cache for storing slices of RR structs.
func NewCache() *Cache {
	cache := &Cache{
		items: make(map[string][]RR),
	}
	go cache.cleanup() // Start cleanup routine to remove expired items.
	return cache
}

// Set adds a slice of RR items to the cache for a specific key. Each RR's Expiry is set based on its TTL.
func (c *Cache) Set(key string, rrs []RR) {
	now := time.Now()
	for i := range rrs {
		if rrs[i].TTL == 0 {
			rrs[i].TTL = time.Second * 86400 * 365 * 100
			rrs[i].Expiry = now.Add(time.Second * 86400 * 365 * 100)
		} else {
			rrs[i].Expiry = now.Add(rrs[i].TTL)
		}
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.items[key] = rrs
}

// Get retrieves a slice of RR items by key if they exist and are unexpired.
func (c *Cache) Get(key string) ([]RR, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	items, found := c.items[key]
	if !found {
		return nil, false
	}
	// Filter out expired records.
	validItems := []RR{}
	for _, item := range items {
		if time.Now().Before(item.Expiry) {
			validItems = append(validItems, item)
		}
	}
	if len(validItems) == 0 {
		return nil, false
	}
	return validItems, true
}

// Delete removes an item from the cache by key.
func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.items, key)
}

// cleanup removes expired items periodically.
func (c *Cache) cleanup() {
	for {
		time.Sleep(time.Minute) // Cleanup interval.
		c.mutex.Lock()
		for key, items := range c.items {
			validItems := []RR{}
			for _, item := range items {
				if time.Now().Before(item.Expiry) {
					validItems = append(validItems, item)
				}
			}
			if len(validItems) > 0 {
				c.items[key] = validItems
			} else {
				delete(c.items, key)
			}
		}
		c.mutex.Unlock()
	}
}

// DumpJSON returns a JSON representation of the current cache state.
func (c *Cache) DumpJSON() (string, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	data, err := json.Marshal(c.items)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// LoadJSON loads the cache state from a JSON string.
func (c *Cache) LoadJSON(data string) error {
	var items map[string][]RR
	if err := json.Unmarshal([]byte(data), &items); err != nil {
		return err
	}

	// Restore the items and adjust Expiry based on current time.
	now := time.Now()
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for key, rrs := range items {
		for i := range rrs {
			// Recompute Expiry based on how much TTL remains.
			timeRemaining := rrs[i].Expiry.Sub(now)
			if timeRemaining > 0 {
				rrs[i].Expiry = now.Add(timeRemaining)
			} else {
				rrs[i].Expiry = now // Expired items get an immediate expiry.
			}
		}
		c.items[key] = rrs
	}
	return nil
}
