package cache

import (
	"social-network/entity"
)

type UserCache interface {
	BaseCache
	Add(key string, object entity.User)
	Get(key string) (entity.User, bool)
	Delete(key string)
}

func (c *Cache[T]) Add(key string, object T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = object
}

func (c *Cache[T]) Get(key string) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, exists := c.cache[key]

	return value, exists
}

func (c *Cache[T]) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.cache, key)
}
