package redis

type Reply interface {
	ToBytes() []byte
}

type ErrReply interface {
	Reply
	Error() string
}
