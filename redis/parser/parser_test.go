package parser

import (
	"bytes"
	"github.com/jiangh156/godis/interface/redis"
	"github.com/jiangh156/godis/redis/protocol"
	"io"
	"reflect"
	"testing"
)

func TestName(t *testing.T) {

}

func TestParseStream(t *testing.T) {
	testReplies := []redis.Reply{
		protocol.MakeIntReply(1),
		protocol.MakeStatusReply("OK"),
		protocol.MakeErrReply("ERR unknown"),
		protocol.MakeBulkReply([]byte("abc")),
		protocol.MakeBulkReply([]byte("")),
		protocol.MakeBulkReply([]byte("a\r\nnb")), //test binary safe
		protocol.MakeNullBulkReply(),
		protocol.MakeMultiBulkReply([][]byte{
			[]byte("a"),
			[]byte("\r\n"),
		}),
		protocol.MakeEmptyMultiBulkReply(),
	}
	testOthers := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		//{
		//	name:     "空行",
		//	input:    []byte("\n"),
		//	expected: nil,
		//},
		//{
		//	name:     "不合理的数字",
		//	input:    []byte(":12s\r\n"),
		//	expected: []byte("illegal number: 12s"),
		//},
		//{
		//	name:     "字符串不存在",
		//	input:    []byte("$-2\r\n"), // $-1\r\n 表示字符串不存在,小于-1的情况
		//	expected: []byte("illegal bulk string header: $-2"),
		//},
		//{
		//	name:     "空字符串数组",
		//	input:    []byte("*-1\r\n"), // *0\r\n 表示空字符串数组,小于0的情况
		//	expected: []byte("illegal array header: *-1"),
		//},
		{
			name:     "test text protocol: set key value",
			input:    []byte("set key value" + protocol.CRLF), // test text protocol: set key value
			expected: protocol.MakeMultiBulkReply([][]byte{[]byte("set"), []byte("key"), []byte("value")}).ToBytes(),
		},
	}
	// 输入reader
	inputReader := bytes.Buffer{}
	for _, v := range testReplies {
		inputReader.Write(v.ToBytes())
	}
	for _, v := range testOthers {
		inputReader.Write(v.input)
	}
	// 期待值
	//expected := make([][]byte, len(testReplies)+len(testOthers))
	expected := [][]byte{}
	for _, v := range testReplies {
		expected = append(expected, v.ToBytes())
	}
	for _, v := range testOthers {
		expected = append(expected, v.expected)
	}

	ch := ParseStream(&inputReader)
	i := 0
	for payload := range ch {
		if payload.Err != nil {
			if payload.Err == io.EOF {
				return
			}
			t.Errorf("test%d occur err, expected:%s, err: %s", i, expected[i], payload.Err)
			//t.Error(payload.Err)
			return
		}
		if !reflect.DeepEqual(payload.Data.ToBytes(), expected[i]) {
			t.Errorf("test%d occur err, expected:%s, got:%s", i, expected[i], payload.Data.ToBytes())
		}
		i++
	}
}
