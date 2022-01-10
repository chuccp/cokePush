package main

import (
	log "github.com/chuccp/coke-log"
	"github.com/chuccp/cokePush/api"
	"github.com/chuccp/cokePush/cluster"
	"github.com/chuccp/cokePush/config"
	"github.com/chuccp/cokePush/core"
	"github.com/chuccp/cokePush/ex"
	"github.com/magiconair/properties"
)

func DefaultRegister() *core.Register {
	cfg, err := config.LoadFile("application.properties", properties.UTF8)
	if err == nil {
		var defaultRegister = core.NewRegister(cfg)
		defaultRegister.AddServer(api.NewServer())
		defaultRegister.AddServer(cluster.NewServer())
		defaultRegister.AddServer(ex.NewServer())
		return defaultRegister
	} else {
		log.PanicF("加载配置文件失败：{}", err.Error())
		return nil
	}

}
func main() {
	config:=log.GetConfig()
	config.SetLevel(log.InfoLevel)
	reg := DefaultRegister()
	cp := reg.Create()
	cp.StartSync()
}
