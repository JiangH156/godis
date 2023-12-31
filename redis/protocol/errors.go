package protocol

// UnKnownErr
type UnknownErrReply struct {
}

var UnknownErrBytes = []byte("-ERR unknown\r\n")

func MakeUnknownErrReply() *UnknownErrReply {
	return &UnknownErrReply{}
}

func (u *UnknownErrReply) ToBytes() []byte {
	return UnknownErrBytes
}

func (u *UnknownErrReply) Error() string {
	return "ERR unknown"
}

// ArgNumErr
type ArgNumErrReply struct {
	Cmd string
}

func MakeArgNumErrReply(cmd string) *ArgNumErrReply {
	return &ArgNumErrReply{Cmd: cmd}
}

func (a *ArgNumErrReply) ToBytes() []byte {
	return []byte("-ERR wrong number of argument for " + a.Cmd + " command\r\n")
}

func (a *ArgNumErrReply) Error() string {
	return "wrong number of argument for " + a.Cmd + " command"
}

// SyntaxErr
type SyntaxErrReply struct {
}

func MakeSyntaxErrReply() *SyntaxErrReply {
	return &SyntaxErrReply{}
}

var SyntaxErrBytes = []byte("-ERR syntax error\r\n")

func (s *SyntaxErrReply) ToBytes() []byte {
	return SyntaxErrBytes
}

func (s *SyntaxErrReply) Error() string {
	return "syntax error"
}

// WrongTypeErr
type WrongTypeErrReply struct {
}

var WrongTypeErrBytes = []byte("-WRONGTYPE Operation against a key holding the wrong king of value\r\n")

func MakeWrongTypeErrReply() *WrongTypeErrReply {
	return &WrongTypeErrReply{}
}
func (w *WrongTypeErrReply) ToBytes() []byte {
	return WrongTypeErrBytes
}

func (w *WrongTypeErrReply) Error() string {
	return "WRONGTYPE Operation against a key holding the wrong king of value"
}

// ProtocolErr
type ProtocolErrReply struct {
	Msg string
}

func MakeProtocolErrReply() *ProtocolErrReply {
	return &ProtocolErrReply{}
}

func (p *ProtocolErrReply) ToBytes() []byte {
	return []byte("-ERR protocol error: '" + p.Msg + "'\r\n")
}

func (p *ProtocolErrReply) Error() string {
	return "ERR protocol error: '" + p.Msg + "'\r\n"
}
