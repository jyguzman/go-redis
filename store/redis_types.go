package store

const (
	STRING int = iota
	LIST
	HASH
	SET
	NIL
)

type RedisValue interface {
	Type() int
}

type RedisString struct {
	Value string
}

func (s RedisString) Type() int { return STRING }

type RedisList struct {
	Values []RedisString
}

func (l *RedisList) Type() int { return LIST }

func (l *RedisList) Push(rv RedisString) {
	l.Values = append(l.Values, rv)
}

func (l *RedisList) Prepend(rv RedisString) {
	l.Values = append([]RedisString{rv}, l.Values...)
}
