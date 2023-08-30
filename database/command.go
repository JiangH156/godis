package database

import "strings"

// 维护redis功能集合，key统一为小写
var cmdTable = make(map[string]*command)
var singleCommand = make(map[string]struct{})

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

func SingleCommandExist(cmdName string) bool {
	cmdName = strings.ToLower(cmdName)
	_, ok := singleCommand[cmdName]
	return ok
}

// PING FLUSHDB FLUSHALL
func RegisterSingleCommand(cmdName string) {
	cmdName = strings.ToLower(cmdName)
	singleCommand[cmdName] = struct{}{}
}
