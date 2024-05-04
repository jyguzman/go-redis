package data_types

type RedisHash struct {
	pairs map[string]RedisValue
}

func NewRedisHash() *RedisHash {
	return &RedisHash{pairs: make(map[string]RedisValue)}
}

func (h *RedisHash) Type() int { return HASH }

func (h *RedisHash) Pairs() map[string]RedisValue {
	return h.pairs
}

func (h *RedisHash) Keys() []string {
	keys := make([]string, len(h.pairs))
	for k := range h.pairs {
		keys = append(keys, k)
	}
	return keys
}

func (h *RedisHash) Set(key string, value RedisString) {
	h.pairs[key] = value
}

func (h *RedisHash) Get(key string) (RedisValue, bool) {
	val, ok := h.pairs[key]
	return val, ok
}

func (h *RedisHash) Delete(key string) int {
	if _, ok := h.Get(key); ok {
		delete(h.pairs, key)
		return 1
	}
	return 0
}
