package core

import (
	log "github.com/chuccp/coke-log"
	"github.com/chuccp/cokePush/config"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Register struct {
	servers   *sync.Map
	context  *Context
}

func (register *Register) AddServer(server Server) {
	server.Init(register.context)
	register.servers.LoadOrStore(server.Name(), server)
}
func (register *Register) Create() *CokePush {
	return &CokePush{register: register,context:register.context }
}
func NewRegister(config *config.Config) *Register {
	return &Register{servers: new(sync.Map),context: newContext(config)}
}

type CokePush struct {
	register *Register
	context  *Context
}

func (cokePush *CokePush) Start() {
	cokePush.context.Init()
	cokePush.register.servers.Range(func(key, value interface{}) bool {
		server, ok := value.(Server)
		if ok {
			var err error
			go func() {
				log.InfoF("启动 {} 服务",server.Name())
				err=server.Start()
				if err!=nil{
					log.ErrorF("启动 {} 服务失败 {}",server.Name(),err.Error())
				}
			}()
		}
		return ok
	})

}
func (cokePush *CokePush) StartSync() {
	cokePush.Start()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGBUS)
	<-sig
}
