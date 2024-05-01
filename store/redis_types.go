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
	Head *RedisListNode
	Tail *RedisListNode
	Size int
}

type RedisListNode struct {
	Value RedisString
	Next  *RedisListNode
	Prev  *RedisListNode
}

func (l *RedisList) Type() int { return LIST }

func (l *RedisList) Append(rs RedisString) int {
	if l.Size == 0 {
		l.Head = &RedisListNode{Value: rs, Next: nil}
		l.Tail = l.Head
		l.Size += 1
		return l.Size
	}

	newNode := &RedisListNode{Value: rs, Next: nil}
	newNode.Prev = l.Tail
	l.Tail.Next = newNode
	l.Tail = newNode
	l.Size += 1
	return l.Size
}

func (l *RedisList) Prepend(rs RedisString) int {
	if l.Size == 0 {
		l.Head = &RedisListNode{Value: rs, Next: nil}
		l.Tail = l.Head
		l.Size += 1
		return 1
	}

	newNode := &RedisListNode{Value: rs, Next: nil}
	newNode.Next = l.Head
	l.Head.Prev = newNode
	l.Head = newNode
	l.Size += 1
	return l.Size
}

//func (l *RedisList) Push(rv RedisString) {
//	l.Values = append(l.Values, rv)
//}
//
//func (l *RedisList) Prepend(rv RedisString) {
//	l.Values = append([]RedisString{rv}, l.Values...)
//}
