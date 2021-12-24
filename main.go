package main

import (
	"github.com/chuccp/cokePush/api"
	"github.com/chuccp/cokePush/config"
	"github.com/chuccp/cokePush/core"
	clog "github.com/chuccp/cokePush/log"
)

func DefaultRegister() *core.Register {
	config := config.DefaultConfig()
	var defaultRegister = core.NewRegister()
	defaultRegister.AddServer(api.NewServer(config))
	return defaultRegister
}
func main() {
	clog.Start()
	reg:=DefaultRegister()
	cp:=reg.Create()
	cp.StartSync()
}
