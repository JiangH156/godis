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
			err := parseBulkString(reader, ch, line)
			if err != nil {
				// 不属于协议错误,直接返回
				ch <- &PayLoad{Err: err}
				return
			}
		case '*': //数组
			err = parseArray(reader, ch, line)
			if err != nil {
				// 不属于协议错误,直接返回
				ch <- &PayLoad{Err: err}
				return
			}
		default:
			args := bytes.Split(line, []byte(" "))
			ch <- &PayLoad{
				Data: protocol.MakeMultiBulkReply(args),
			}
		}
	}

}

// 读取BulkString
func parseBulkString(reader *bufio.Reader, ch chan<- *PayLoad, header []byte) error {
	msgLen, err := strconv.ParseInt(string(header[1:]), 10, 64)
	// $-1\r\n表示一个不存在的值
	// $0\r\n表示一个空字符串。
	if err != nil || msgLen < -1 {
		protocolError(ch, "illegal bulk string header: "+string(header))
		return nil
	} else if msgLen == -1 {
		ch <- &PayLoad{
			Data: protocol.MakeNullBulkReply(), //$-1\r\n表示一个不存在的值
		}
		return nil
	} else {
		msg := make([]byte, int(msgLen)+2)
		_, err = io.ReadFull(reader, msg)
		if err != nil {
			return err
		}
		ch <- &PayLoad{
			Data: protocol.MakeBulkReply(msg[:len(msg)-2]),
		}
	}
	return nil
}
func parseArray(reader *bufio.Reader, ch chan<- *PayLoad, header []byte) error {
	arrLen, err := strconv.ParseInt(string(header[1:]), 10, 64)
	//fmt.Printf("%s\n", header)
	if err != nil || arrLen < 0 {
		protocolError(ch, "illegal array header: "+string(header))
		return nil
	}
	// *0\r\n 表示空数组
	if arrLen == 0 {
		ch <- &PayLoad{Data: protocol.MakeEmptyMultiBulkReply()}
		return nil
	}
	args := [][]byte{}
	for i := int64(0); i < arrLen; i++ {
		var line []byte
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return err
		}
		length := len(line)
		// 考虑空行的情况
		if length < 4 || line[length-2] != '\r' || line[0] != '$' {
			protocolError(ch, "illegal bulk string header: "+string(line))
		}
		// 获取字符串长度
		strLen, err := strconv.ParseInt(string(line[1:length-2]), 10, 64)
		if err != nil || strLen < -1 {
			protocolError(ch, "illegal bulk string length: "+string(line))
		} else if strLen == -1 { //字符串不存在,为空,且后续没有字符串数据,不需要再次读取
			args = append(args, []byte{})
		} else {
			arg := make([]byte, strLen+2)
			_, err = io.ReadFull(reader, arg)
			if err != nil {
				return err
			}
			if arg[len(arg)-2] != '\r' || arg[len(arg)-1] != '\n' {
				protocolError(ch, "illegal bulk string: "+string(arg))
				return nil
			}
			args = append(args, arg[:len(arg)-2])
		}
	}
	ch <- &PayLoad{Data: protocol.MakeMultiBulkReply(args)}
	return nil
}

func protocolError(ch chan<- *PayLoad, msg string) {
	err := errors.New("protocol error: " + msg)
	ch <- &PayLoad{Err: err}
}
