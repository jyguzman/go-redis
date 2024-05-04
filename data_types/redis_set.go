package data_types

type RedisSet struct {
	values map[string]RedisString
	size   int
}

func (set *RedisSet) Type() int { return SET }

func (set *RedisSet) Size() int { return set.size }

func (set *RedisSet) Values() map[string]RedisString {
	return set.values
}

func NewRedisSet() *RedisSet {
	return &RedisSet{values: make(map[string]RedisString), size: 0}
}

func (set *RedisSet) Add(value string) int {
	if set.Contains(value) {
		return 0
	}
	set.values[value] = NewRedisString(value)
	set.size++
	return 1
}

func (set *RedisSet) Remove(value string) int {
	if !set.Contains(value) {
		return 0
	}
	delete(set.values, value)
	set.size--
	return 1
}

func (set *RedisSet) Contains(value string) bool {
	_, ok := set.values[value]
	return ok
}
