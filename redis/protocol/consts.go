package protocol

// PING
type PingReply struct {
}

var pingBytes = []byte("+PONG\r\n")

func (p *PingReply) ToBytes() []byte {
	return pingBytes
}

// OK
type OkReply struct {
}

var okBytes = []byte("+OK\r\n")

func (o *OkReply) ToBytes() []byte {
	return okBytes
}

// NullBulk
type NullBulkReply struct {
}

var nullBulkBytes = []byte("$-1\r\n")

func (n *NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

// EmptyMultiBulk
type EmptyMultiBulkReply struct {
}

var emptyMultiBulkBytes = []byte("*0\r\n")

func (e *EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}

//	NoReply
//
// reply nothing, for commands like subscribe
type NoReply struct {
}

var noBytes = []byte("")

func (n *NoReply) ToBytes() []byte {
	return noBytes
}
