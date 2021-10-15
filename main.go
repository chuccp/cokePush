package main

import (
	"github.com/chuccp/cokePush/api"
	"github.com/chuccp/cokePush/config"
	"github.com/chuccp/cokePush/core"
)



func DefaultRegister() *core.Register {
	config:=config.DefaultConfig()
	var defaultRegister = core.NewRegister()
	defaultRegister.AddServer(api.NewServer(config))
	return defaultRegister
}
func main() {
	reg:=DefaultRegister()
	cokePush:=reg.Create()
	cokePush.StartSync()
}
