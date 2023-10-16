package main

import (
	"fmt"
	"github.com/jiangh156/godis/config"
	"github.com/jiangh156/godis/lib/logger"
	"github.com/jiangh156/godis/redis/server"
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
	logger.Setup(&logger.LogCfg{
		Path:       "logs",
		Name:       "godis",
		Ext:        ".log",
		TimeFormat: "2006-01-01",
	})
	config.SetupConfig("redis.yml")
	//handler := server.MakeRedisHandler()
	handler := server.MakeRedisHandler()
	err := tcp.ListenAndServeWithSignal(&tcp.Config{
		Address: config.Properties.Bind,
		Port:    config.Properties.Port,
		MaxConn: 10,
	}, handler)
	if err != nil {
		fmt.Println(err)
	}
}
