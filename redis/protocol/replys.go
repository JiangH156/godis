package protocol

import (
	"strconv"
)

/*
简单字符串：以"+" 开始， 如："+OK\r\n"
错误：以"-" 开始，如："-ERR Invalid Synatx\r\n"
整数：以":"开始，如：":1\r\n"
字符串：以 $ 开始
数组：以 * 开始
*/
var (
	nullBulkReplyBytes = []byte("$-1\r\n")
	CRLF               = "\r\n"
)

// Bulk Reply
type BulkReply struct {
	Arg []byte
}

func MakeBulkReply(arg []byte) *BulkReply {
	return &BulkReply{
		Arg: arg,
	}
}
func (b *BulkReply) ToBytes() []byte {
	if len(b.Arg) == 0 {
		return nullBulkReplyBytes
	}
	return []byte("$" + strconv.Itoa(len(b.Arg)) + CRLF + string(b.Arg) + CRLF)
}

// multiBulk reply
type MultiBulkReply struct {
	Args [][]byte
}

func (m *MultiBulkReply) ToBytes() []byte {
	res := "*" + strconv.Itoa(len(m.Args)) + CRLF
	for _, arg := range m.Args {
		// arg为空
		if len(arg) == 0 {
			res += "$-1" + CRLF
		} else {
			res += "$" + strconv.Itoa(len(arg)) + CRLF + string(arg) + CRLF
		}
	}
	return []byte(res)
}

func MakeMultiBulkReply(args [][]byte) *MultiBulkReply {
	return &MultiBulkReply{
		Args: args,
	}
}

// int Reply
type IntReply struct {
	Code int64
}

func MakeIntReply(code int64) *IntReply {
	return &IntReply{
		Code: code,
	}
}

func (i *IntReply) ToBytes() []byte {
	return []byte(":" + strconv.FormatInt(i.Code, 10) + CRLF)
}

// status Reply
type StatusReply struct {
	Status string
}

func (s *StatusReply) ToBytes() []byte {
	return []byte("$" + s.Status + CRLF)
}

func MakeStatusReply(status string) *StatusReply {
	return &StatusReply{
		Status: status,
	}
}

// error Reply
type StandardErrorReply struct {
	Status string
}

func MakeErrReply(status string) *StandardErrorReply {
	return &StandardErrorReply{
		Status: status,
	}
}

func (s *StandardErrorReply) ToBytes() []byte {
	return []byte("-" + s.Status + CRLF)
}

func (s *StandardErrorReply) Error() string {
	return s.Status
}
