package store

type Cache struct {
	store map[string][]byte
}

func (c *Cache) Add(key string, value []byte) {}

func (c *Cache) Contains(key string) bool {
	return c.store[key] != nil
}

func (c *Cache) Get(key string) {}

func (c *Cache) Remove(key string) {}
