package db

import "strings"

// 维护redis功能集合，key统一为小写
var cmdTable = make(map[string]*command)

type command struct {
	exector ExecFunc
	arity   int
}

func RegisterCommand(name string, exector ExecFunc, arity int) {
	name = strings.ToLower(name)
	cmdTable[name] = &command{
		exector: exector,
		arity:   arity,
	}
}
