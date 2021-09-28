package main

import (
	"github.com/chuccp/cokePush/config"
	"github.com/chuccp/cokePush/server"
)

func main() {

	config:=config.DefaultConfig()
	server.Start(config)
}
