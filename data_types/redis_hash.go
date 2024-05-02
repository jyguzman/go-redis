package data_types

type RedisHash struct {
	Pairs map[string]RedisValue
}

func NewRedisHash() *RedisHash {
	return &RedisHash{Pairs: make(map[string]RedisValue)}
}

func (h *RedisHash) Type() int { return HASH }

func (h *RedisHash) Set(key string, value RedisString) {
	h.Pairs[key] = value
}

func (h *RedisHash) Get(key string) (RedisValue, bool) {
	val, ok := h.Pairs[key]
	return val, ok
}
