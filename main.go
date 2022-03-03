package main

import (
	"github.com/chuccp/cokePush/api"
	"github.com/chuccp/cokePush/cluster"
	"github.com/chuccp/cokePush/core"
	"github.com/chuccp/cokePush/ex"
	"github.com/chuccp/cokePush/tcp"
	"github.com/chuccp/utils/config"
	"github.com/chuccp/utils/log"
	"github.com/magiconair/properties"
)

func DefaultRegister() *core.Register {
	cfg, err := config.LoadFile("application.properties", properties.UTF8)
	if err == nil {
		var defaultRegister = core.NewRegister(cfg)
		defaultRegister.AddServer(api.NewServer())
		defaultRegister.AddServer(tcp.NewServer())
		defaultRegister.AddServer(cluster.NewServer())
		defaultRegister.AddServer(ex.NewServer())
		return defaultRegister
	} else {
		log.PanicF("加载配置文件失败：{}", err.Error())
		return nil
	}

}
func main() {
	config:=log.GetDefaultConfig()
	config.SetLevel(log.InfoLevel)
	reg := DefaultRegister()
	cp := reg.Create()
	cp.StartSync()
}
