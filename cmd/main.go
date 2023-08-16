package main

import (
	"fmt"
	"github.com/jiangh156/godis/config"
	"github.com/jiangh156/godis/tcp"
)

func main() {
	config.SetupConfig("redis.yml")
	handler := &tcp.EchoHandler{}
	err := tcp.ListenAndServeWithSignal(&tcp.Config{
		Address: config.Properties.Bind,
		Port:    config.Properties.Port,
		MaxConn: 10,
	}, handler)
	if err != nil {
		fmt.Println(err)
	}
}
