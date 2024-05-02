package data_types

type RedisSet struct {
	Values map[string]RedisString
	Size   int
}

func (set *RedisSet) Type() int {
	return SET
}

func NewRedisSet() *RedisSet {
	return &RedisSet{Values: make(map[string]RedisString), Size: 0}
}

func (set *RedisSet) Add(value string) int {
	if set.Contains(value) {
		return 0
	}
	set.Values[value] = NewRedisString(value)
	set.Size++
	return 1
}

func (set *RedisSet) Remove(value string) int {
	if !set.Contains(value) {
		return 0
	}
	delete(set.Values, value)
	set.Size--
	return 1
}

func (set *RedisSet) Contains(value string) bool {
	_, ok := set.Values[value]
	return ok
}
