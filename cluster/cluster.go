package cluster

import (
	log "github.com/chuccp/coke-log"
	"github.com/chuccp/cokePush/config"
	"github.com/chuccp/cokePush/core"
	"github.com/chuccp/cokePush/km"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/net"
	"github.com/chuccp/cokePush/user"
	"strconv"
	"strings"
	"time"
)

func newMachine(remotePort int, remoteHost string) *machine {
	return &machine{remotePort: remotePort, remoteHost: remoteHost, isLocal: false}
}

type Server struct {
	config     *config.Config
	port       int
	machineId  string
	tcpserver  *net.TCPServer
	userStore  *user.Store
	machine    *machine
	machineMap *machineStore
	request    *km.Request
	context *core.Context
}
type MachineInfo struct {
	Address string
	UserNum uint
}

func (server *Server) Start() error {
	err := server.tcpserver.Bind()
	if err != nil {
		return err
	}
	go server.AcceptConn()
	go server.queryMachine()
	return nil
}

func (server *Server) AcceptConn() {
	for {
		io, err := server.tcpserver.Accept()
		if err != nil {
			break
		} else {
			client, err := NewClient(io, server)
			if err == nil {
				go client.Start()
			}
		}
	}
}
func (server *Server)machineInfo(value ...interface{})interface{}   {

	mis:=make([]MachineInfo,0)

	server.machineMap.each(func(machineId string, machine *machine) bool {
		var mi MachineInfo
		mi.Address = machine.remoteHost+":"+strconv.Itoa(machine.remotePort)+"|"+machineId

		log.TraceF("请求 host:{}，port:{}",machine.remoteHost,machine.remotePort)
		server.queryMachineBasic1(machine.remoteHost,machine.remotePort)

		mis = append(mis, mi)
		return true
	})
	var mi MachineInfo
	mi.Address = "localhost"+":"+strconv.Itoa(server.port)+"|"+server.machineId
	mi.UserNum = server.userStore.GetUserNum()
	mis = append(mis, mi)
	return mis
}

func (server *Server) queryMachine() {
	var hasQuery = false
	for {
		time.Sleep(time.Second * 2)
		log.DebugF("查询机器信息 hasQuery:{} ", hasQuery)
		if !hasQuery && !server.machine.isLocal {
			err := server.queryMachineBasic()
			log.DebugF("查询机器信息 hasQuery:{} ", hasQuery, err)
			if err == nil {
				hasQuery = true
			} else {
				log.ErrorF("queryMachineInfo err:{}", err)
			}
		}
		if hasQuery {
			server.getMachineList()
		}
		time.Sleep(time.Minute)
	}
}
func (server *Server) queryMachineBasic1(host string,port int)(*machine,*km.Conn,error){
	qBasic := newQueryMachineBasic(server.port, server.machineId)
	log.DebugF("queryMachineInfo发送信息 :{} msgId:{}", server.machineId, qBasic.GetMessageId())
	msg, conn, err := server.request.Call(host, port, qBasic)
	log.DebugF("queryMachineInfo 收到信息 :{}", server.machineId)
	if err==nil{
		machine := msg.GetString(message.BackMachineAddress)
		m, err := toMachine(machine)
		if err==nil{
			return m,conn,nil
		}else{
			return nil,nil,err
		}
	}
	return nil,nil,err
}
/**客户端获取请求机器信息**/
func (server *Server) queryMachineBasic() error {

	qBasic := newQueryMachineBasic(server.port, server.machineId)
	log.DebugF("queryMachineInfo发送信息 :{} msgId:{}", server.machineId, qBasic.GetMessageId())
	msg, conn, err := server.request.Call(server.machine.remoteHost, server.machine.remotePort, qBasic)
	log.DebugF("queryMachineInfo 收到信息 :{}", server.machineId)
	if err == nil {
		if message.BackMessageClass == msg.GetClassId() {
			if message.BackMessageOKType == msg.GetMessageType() {
				machine := msg.GetString(message.BackMachineAddress)
				if len(machine) > 0 {
					m, err := toMachine(machine)
					if err == nil {
						if m.machineId == server.machineId {
							m.isLocal = true
							server.machine.isLocal = true
							log.InfoF("连接到自己关闭连接 :{}", machine)
							conn.Close()
							return nil
						} else {
							m.remoteHost = server.machine.remoteHost
							m.remotePort = server.machine.remotePort
							if server.machineMap.add(m) {
								log.InfoF("queryMachineInfo 添加新的机器连接", machine)
							}
							return nil
						}
					} else {
						return err
					}
				}
			}
		}
	} else {
		return err
	}
	return core.UnKnownConn
}

/**客户端获取机器列表**/
func (server *Server) getMachineList() {
	qMsg := newQueryMachineMessage(server.port, server.machineId)
	server.machineMap.eachAddress(func(remoteHost string, remotePort int) {
		log.InfoF("!!!!!!!getMachineList  remoteHost：{} remotePort：{}", remoteHost,remotePort)
		msg, _, err := server.request.Call(remoteHost, remotePort, qMsg)
		if err != nil {
			log.ErrorF("getMachineList err:{}  remoteHost：{} remotePort：{}", err,remoteHost,remotePort)
		} else {
			if message.BackMessageClass == msg.GetClassId() {
				if message.BackMessageOKType == msg.GetMessageType() {
					machines := msg.GetString(message.BackMachineAddress)
					if len(machines) > 0 {
						addresses := strings.Split(machines, ";")
						if len(addresses) > 0 {
							for _, v := range addresses {
								m, err := toMachine(v)
								if err != nil {
									log.ErrorF("getMachineList 解析地址错误 GetMessageType:{}  err:{}", v, err)
								} else {
									if m.machineId == server.machineId {
										m.isLocal = true
									} else {
										if server.machineMap.add(m) {
											log.InfoF("getMachineList 添加新的机器连接", v)
										}
									}
								}
							}
						}
					}
				} else {
					log.ErrorF("getMachineList 未知错误 GetMessageType：{}", msg.GetMessageType())
				}
			} else {
				log.ErrorF("getMachineList 未知返回 BackMessageClass：{}", msg.GetClassId())
			}
		}
	})
}

// Init 初始化
func (server *Server) Init(context *core.Context) {
	server.context = context
	server.port = server.config.GetIntOrDefault("cluster.local.port", 6361)
	server.machineId = server.config.GetStringOrDefault("cluster.local.machineId", MachineId())
	if server.machineId == "" {
		log.PanicF("machineId can‘t blank")
	}
	var remotePort = server.config.GetIntOrDefault("cluster.remote.port", 6362)
	var remoteHost = server.config.GetString("cluster.remote.host")
	if remoteHost == "" {
		log.PanicF("cluster.remote.host can‘t blank")
	}
	server.request = km.NewRequest()
	server.machineMap = newMachineStore()
	server.machine = newMachine(remotePort, remoteHost)
	server.tcpserver = net.NewTCPServer(server.port)
	context.RegisterHandle("machineInfo",server.machineInfo)
}
func (server *Server) Name() string {
	return "cluster"
}

func NewServer(config *config.Config) *Server {
	return &Server{config: config}
}
