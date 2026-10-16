package cache

import (
	"maps"
	"social-network/entity"
	"sync"
)

type BaseCache interface {
	Len() int
	GetData() map[string]entity.User
}

type Cache[T any] struct {
	mu    *sync.RWMutex
	cache map[string]T
}

func NewCache[T any]() *Cache[T] {
	return &Cache[T]{
		cache: make(map[string]T),
		mu:    &sync.RWMutex{},
	}
}

func (c *Cache[T]) Len() int {
	return len(c.cache)
}

func (c *Cache[T]) GetData() map[string]T {
	c.mu.RLock()
	defer c.mu.RUnlock()

	data := make(map[string]T, len(c.cache))
	maps.Copy(data, c.cache)

	return data
}
