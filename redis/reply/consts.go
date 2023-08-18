package reply

// PING
type PingReply struct {
}

var PingBytes = []byte("+PONG\r\n")

func (p *PingReply) ToBytes() []byte {
	return PingBytes
}

// OK
type OkReply struct {
}

var OKBytes = []byte("+OK\r\n")

func (o *OkReply) ToBytes() []byte {
	return OKBytes
}

// NullBulk
type NullBulkReply struct {
}

var NullBulkBytes = []byte("$1\r\n")

func (n *NullBulkReply) ToBytes() []byte {
	return NullBulkBytes
}

// EmptyMultiBulk
type EmptyMultiBulkReply struct {
}

var EmptyMultiBulkBytes = []byte("*0\r\n")

func (e *EmptyMultiBulkReply) ToBytes() []byte {
	return EmptyMultiBulkBytes
}

//	NoReply
//
// reply nothing, for commands like subscribe
type NoReply struct {
}

var NoBytes = []byte("")

func (n *NoReply) ToBytes() []byte {
	return NoBytes
}
