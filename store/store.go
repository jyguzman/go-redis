package store

type Cache struct {
	store map[string]RedisValue
}

func (c *Cache) Contains(key string) bool {
	return c.store[key] != nil
}

func (c *Cache) Set(key string, value RedisValue) { c.store[key] = value }

func (c *Cache) Get(key string) RedisValue { return c.store[key] }

func (c *Cache) Remove(key string) { delete(c.store, key) }

var Store = Cache{store: make(map[string]RedisValue)}
