package core

import "sync"

type Register struct {
	servers *sync.Map
	apiServer Server
}

func (register *Register) AddServer(server Server) {
	if server.Name()=="api"{
		register.apiServer = server
	}else{
		register.servers.LoadOrStore(server.Name(),server)
	}
}
func (register *Register) Create()*CokePush {
	return  &CokePush{register:register}
}
func NewRegister()*Register  {

	return &Register{servers: new(sync.Map)}
}

type CokePush struct {
	register *Register
}
func (cokePush *CokePush) Start() {
	cokePush.register.servers.Range(func(key, value interface{}) bool {
		server,ok:=value.(Server)
		if ok{
			server.Init()
			go server.Start()
		}
		return ok
	})
	if cokePush.register.apiServer!=nil{
		cokePush.register.apiServer.Init()
		go cokePush.register.apiServer.Start()
	}
}
func (cokePush *CokePush) StartSync()  {
	cokePush.register.servers.Range(func(key, value interface{}) bool {
		server,ok:=value.(Server)
		if ok{
			server.Init()
			go server.Start()
		}
		return ok
	})
	if cokePush.register.apiServer!=nil{
		cokePush.register.apiServer.Init()
		cokePush.register.apiServer.Start()
	}
}