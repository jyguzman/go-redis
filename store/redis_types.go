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
	RedisString RedisString
	Next        *RedisListNode
	Prev        *RedisListNode
}

func (l *RedisList) Type() int { return LIST }

func (l *RedisList) Append(rs RedisString) int {
	if l.Size == 0 {
		l.Head = &RedisListNode{RedisString: rs, Next: nil}
		l.Tail = l.Head
		l.Size += 1
		return l.Size
	}

	newNode := &RedisListNode{RedisString: rs, Next: nil}
	newNode.Prev = l.Tail
	l.Tail.Next = newNode
	l.Tail = newNode
	l.Size += 1
	return l.Size
}

func (l *RedisList) Prepend(rs RedisString) int {
	if l.Size == 0 {
		l.Head = &RedisListNode{RedisString: rs, Next: nil}
		l.Tail = l.Head
		l.Size += 1
		return 1
	}

	newNode := &RedisListNode{RedisString: rs, Next: nil}
	newNode.Next = l.Head
	l.Head.Prev = newNode
	l.Head = newNode
	l.Size += 1
	return l.Size
}

func (l *RedisList) PopLeft() string {
	val := l.Head.RedisString.Value
	l.Head = l.Head.Next
	return val
}

func (l *RedisList) Pop() string {
	val := l.Tail.RedisString.Value
	l.Tail = l.Tail.Prev
	return val
}
