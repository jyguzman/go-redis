package data_types

const (
	STRING int = iota
	LIST
	HASH
	SET
)

type RedisValue interface {
	Type() int
}

func CheckType(val RedisValue) bool {
	return false
}

type RedisString struct {
	Value string
}

func NewRedisString(str string) RedisString {
	return RedisString{Value: str}
}

func (rs RedisString) Type() int { return STRING }
