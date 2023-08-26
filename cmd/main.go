package main

import (
	"fmt"

	"github.com/jiangh156/godis/config"
	"github.com/jiangh156/godis/tcp"
)

var banner = `
   ______          ___
  / ____/___  ____/ (_)____
 / / __/ __ \/ __  / / ___/
/ /_/ / /_/ / /_/ / (__  )
\____/\____/\__,_/_/____/
`

func main() {
	print(banner)
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
