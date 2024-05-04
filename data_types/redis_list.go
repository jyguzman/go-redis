package data_types

type RedisList struct {
	Head *RedisListNode
	Tail *RedisListNode
	size int
}

type RedisListNode struct {
	RedisString RedisString
	Next        *RedisListNode
	Prev        *RedisListNode
}

func NewRedisList() *RedisList {
	return &RedisList{size: 0}
}

func (l *RedisList) Type() int { return LIST }

func (l *RedisList) Size() int { return l.size }

func (l *RedisList) Append(rs RedisString) int {
	if l.size == 0 {
		l.Head = &RedisListNode{RedisString: rs, Next: nil}
		l.Tail = l.Head
		l.size += 1
		return l.size
	}

	newNode := &RedisListNode{RedisString: rs, Next: nil}
	newNode.Prev = l.Tail
	l.Tail.Next = newNode
	l.Tail = newNode
	l.size += 1
	return l.size
}

func (l *RedisList) Prepend(rs RedisString) int {
	if l.size == 0 {
		l.Head = &RedisListNode{RedisString: rs}
		l.Tail = l.Head
		l.size += 1
		return 1
	}

	newNode := &RedisListNode{RedisString: rs}
	newNode.Next = l.Head
	l.Head.Prev = newNode
	l.Head = newNode
	l.size += 1
	return l.size
}

func (l *RedisList) PopLeft() string {
	val := l.Head.RedisString.Value
	l.Head = l.Head.Next
	l.size--
	return val
}

func (l *RedisList) Pop() string {
	val := l.Tail.RedisString.Value
	l.Tail = l.Tail.Prev
	l.size--
	return val
}
