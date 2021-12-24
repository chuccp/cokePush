package core

import (
	log "github.com/chuccp/coke-log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Register struct {
	servers   *sync.Map
}

func (register *Register) AddServer(server Server) {
	register.servers.LoadOrStore(server.Name(), server)
}
func (register *Register) Create() *CokePush {
	return &CokePush{register: register, context: newContext()}
}
func NewRegister() *Register {

	return &Register{servers: new(sync.Map)}
}

type CokePush struct {
	register *Register
	context  *Context
}

func (cokePush *CokePush) Start() {
	cokePush.register.servers.Range(func(key, value interface{}) bool {
		server, ok := value.(Server)
		if ok {
			server.Init(cokePush.context)
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
