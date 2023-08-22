package protocol

// PING
type PingReply struct {
}

var pingBytes = []byte("+PONG\r\n")

func MakePingReply() *PingReply {
	return &PingReply{}
}

func (p *PingReply) ToBytes() []byte {
	return pingBytes
}

// OK
type OkReply struct {
}

var okBytes = []byte("+OK\r\n")

func MakeOkReply() *OkReply {
	return &OkReply{}
}

func (o *OkReply) ToBytes() []byte {
	return okBytes
}

// NullBulk
type NullBulkReply struct {
}

var nullBulkBytes = []byte("$-1\r\n")

func MakeNullBulkReply() *NullBulkReply {
	return &NullBulkReply{}
}

func (n *NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

// EmptyMultiBulk
type EmptyMultiBulkReply struct {
}

var emptyMultiBulkBytes = []byte("*0\r\n")

func MakeEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return &EmptyMultiBulkReply{}
}

func (e *EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}

//	NoReply
//
// reply nothing, for commands like subscribe
type NoReply struct {
}

var noBytes = []byte("")

func MakeNoReply() *NoReply {
	return &NoReply{}
}

func (n *NoReply) ToBytes() []byte {
	return noBytes
}
