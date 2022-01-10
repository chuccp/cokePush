package main

import (
	log "github.com/chuccp/coke-log"
	"github.com/chuccp/cokePush/api"
	"github.com/chuccp/cokePush/cluster"
	"github.com/chuccp/cokePush/config"
	"github.com/chuccp/cokePush/core"
	"github.com/chuccp/cokePush/ex"
	clog "github.com/chuccp/cokePush/log"
	"github.com/magiconair/properties"
)

func DefaultRegister() *core.Register {
	config,err := config.LoadFile("application.properties",properties.UTF8)
	if err==nil{
		var defaultRegister = core.NewRegister(config)
		defaultRegister.AddServer(api.NewServer())
		defaultRegister.AddServer(cluster.NewServer())
		defaultRegister.AddServer(ex.NewServer())
		return defaultRegister
	}else{
		log.PanicF("加载配置文件失败：{}",err.Error())
		return nil
	}

}
func main() {
	clog.Start()
	reg:=DefaultRegister()
	cp:=reg.Create()
	cp.StartSync()
}
