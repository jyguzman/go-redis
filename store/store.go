package store

import "go-redis/data_types"

type Cache struct {
	store map[string]data_types.RedisValue
}

var Store = &Cache{store: make(map[string]data_types.RedisValue)}

func (c *Cache) Contains(key string) bool {
	return c.store[key] != nil
}

func (c *Cache) Set(key string, value data_types.RedisValue) { c.store[key] = value }

func (c *Cache) Get(key string) data_types.RedisValue { return c.store[key] }

func (c *Cache) Type(key string) int { return c.Get(key).Type() }

func (c *Cache) Remove(key string) { delete(c.store, key) }
