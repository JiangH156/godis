package parser

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/jiangh156/godis/interface/redis"
	"github.com/jiangh156/godis/lib/logger"
	"github.com/jiangh156/godis/redis/protocol"
	"io"
	"runtime/debug"
	"strconv"
)

type PayLoad struct {
	Data redis.Reply
	Err  error
}

//type readState struct {
//	msgType           byte
//	readingMultiLine  bool
//	expectedArgsCount int
//	args              [][]byte
//	bulkLen           int
//}

// ParseStream 通过io.Reader读取数据并将结果通过channel返回给调用者
func ParseStream(reader io.Reader) <-chan *PayLoad {
	ch := make(chan *PayLoad)
	go parse0(reader, ch)
	return ch
}
func parse0(rawReader io.Reader, ch chan<- *PayLoad) {
	go func() {
		if err := recover(); err != nil {
			logger.Error(err, string(debug.Stack()))
		}
	}()
	reader := bufio.NewReader(rawReader)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			ch <- &PayLoad{Err: err}
			close(ch)
			return
		}
		length := len(line)
		// 考虑空行的情况
		if length <= 2 || line[length-2] != '\r' {
			continue
		}
		// 删除后缀
		line = bytes.TrimSuffix(line, []byte{'\r', '\n'})
		switch line[0] {
		case '+': //简单字符串
			ch <- &PayLoad{
				Data: protocol.MakeStatusReply(string(line[1:])), //简单字符串，此处使用StatusReply
			}
		case ':': //整数
			code, err := strconv.ParseInt(string(line[1:]), 10, 64)
			if err != nil {
				protocolError(ch, "illegal number: "+string(line[1:])) // 不合理的数字
				continue
			}
			ch <- &PayLoad{
				Data: protocol.MakeIntReply(code),
			}
		case '-': //错误
			ch <- &PayLoad{
				Data: protocol.MakeErrReply(string(line[1:])),
			}
		case '$': //字符串
		case '*': //数组

		}
	}

}

func parseBulkString(msg []byte) error {
	return nil
}
func parseMultiBulkString(msg []byte) error {
	return nil
}

func protocolError(ch chan<- *PayLoad, msg string) {
	err := errors.New("protocol error: " + msg)
	ch <- &PayLoad{Err: err}
}
