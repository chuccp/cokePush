package cluster

import (
	log "github.com/chuccp/coke-log"
	"github.com/chuccp/cokePush/core"
	"github.com/chuccp/cokePush/km"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/net"
	"github.com/chuccp/cokePush/user"
	"github.com/chuccp/cokePush/util"
	"github.com/pquerna/ffjson/ffjson"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

func newMachine(remotePort int, remoteHost string) *machine {
	return &machine{remotePort: remotePort, remoteHost: remoteHost, isLocal: false}
}

type Server struct {
	port       int
	machineId  string
	tcpserver  *net.TCPServer
	machine    *machine
	machineMap *machineStore
	request    *km.Request
	context    *core.Context
	userStore *userStore
}
type MachineInfo struct {
	Address string
	UserNum int32
	MachineId string
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

func (server *Server) HandleAddUser(iUser user.IUser) {
	log.InfoF("发送添加用户")
	msg:=newAddUserMessage(server.machineId,iUser.GetUsername())
	server.machineMap.eachAddress(func(remoteHost string, remotePort int){
		server.request.JustCall(remoteHost,remotePort,msg)
	})
}
func (server *Server) Query(value ...interface{}) interface{}{
	log.InfoF("集群查询 {}",value)
	queryName:=value[0].(string)
	v:=value[1]
	params :=make([]string,0)
	if len(value)>2{
		vs:=value[2:len(value)]
		for _,vv:=range vs{
			params = append(params,vv.(string) )
		}
	}
	query:=newQuery(queryName, params...)
	vvs:=make([]interface{},0)
	vvs = append(vvs, v)
	server.machineMap.eachAddress(func(remoteHost string, remotePort int) {
		log.InfoF("集群查询 {}:{}",remoteHost,remotePort)
		im,_,err:= server.request.Call(remoteHost,remotePort,query)
		log.InfoF("集群查询 {}：{} err：{}",remoteHost,remotePort,err)
		if err==nil{
			if im.GetMessageType()==message.BackMessageOKType{
				data:=im.GetValue(message.QueryData)
				if len(data)>0{
					var m =util.NewPtr(v)
					err1:=ffjson.Unmarshal(data,m)
					if err1==nil{
						vvs = append(vvs, m)
					}else{
						log.InfoF("type:{},集群查询 err1: {}",reflect.TypeOf(m),err1)
					}
				}
			}
		}
	})
	return vvs
}
func (server *Server) HandleDeleteUser(username string) {
	msg:=newDeleteUserMessage(server.machineId,username)
	server.machineMap.eachAddress(func(remoteHost string, remotePort int){
		server.request.JustCall(remoteHost,remotePort,msg)
	})
}
func (server *Server)addUser(username string,machineId string){
	m,ok:=server.machineMap.getMachine(machineId)
	if ok{
		server.userStore.AddUser(username,m.remoteHost,m.remotePort)
	}
}
func (server *Server)delete(username string,machineId string)  {
	server.userStore.DeleteUser(username)
}
func (server *Server) HandleSendMessage(iMessage *core.DockMessage, writeFunc user.WriteFunc)  {
	server.sendStoreMachineDockMessage(iMessage, func(err error, hasUser bool,host string,port int) {
		if hasUser{
			writeFunc(err,hasUser)
		}else{
			flag:=server.sendAllMachineDockMessage(iMessage,writeFunc,"",0)
			if !flag{
				writeFunc(nil,false)
			}
		}
	})
}
func (server *Server) HandleSendMultiMessage(fromUser string, usernames *[]string, text string, f func(username string, status int)){
	userMap:=make(map[string]*[]string)
	for _,v:=range *usernames{
		cu,ok:=server.userStore.GetUserMachine(v)
		if ok{
			usernameArray := userMap[cu.remoteAddress]
			if usernameArray==nil{
				us := make([]string,0)
				userMap[cu.remoteAddress] = &us
				usernameArray = &us
			}
			*usernameArray = append(*usernameArray, v)
			f(v,1)
		}else{
			f(v,0)
		}
	}
	go server.sendMultiMessage(fromUser,userMap,text)
}
func (server *Server)sendMultiMessage(fromUser string,userMap map[string]*[]string,text string){
	for k,v:=range userMap{
		server.request.JustCall2(k,message.CreateMultiMessage(fromUser,v,text))
	}
}
func (server *Server)sendStoreMachineDockMessage(iMessage *core.DockMessage,f func(err error, hasUser bool,host string,port int)){
	username:=iMessage.GetToUsername()
	u,ok:=server.userStore.GetUserMachine(username)
	if ok{
		server.request.Async(u.remoteHost, u.remotePort, iMessage.InputMessage, func(replayMessage message.IMessage, hasReplay bool, err error) {
			if hasReplay{
				ty:=replayMessage.GetMessageType()
				if ty==message.BackMessageOKType{
					f(nil,true,u.remoteHost, u.remotePort)
				}else{
					log.DebugF("发信息失败 u:{} DeleteUser",username)
					server.userStore.DeleteUser(username)
					f(err,false,u.remoteHost, u.remotePort)
				}
				log.DebugF("发信息失败 u:{} DeleteUser",username)
			}else{
				server.userStore.DeleteUser(username)
				f(err,false,u.remoteHost, u.remotePort)
			}
		})
	}else{
		f(nil,false,"",0)
	}
}


func (server *Server) sendAllMachineDockMessage(iMessage *core.DockMessage, writeFunc user.WriteFunc,exHost string,exPort int)bool  {
	var i int32 = 0
	var flag = false
	var hasMachine = false
	server.machineMap.eachAddress(func(remoteHost string, remotePort int) {
		if remoteHost==exHost &&remotePort==exPort{
			return
		}
		hasMachine = true
		atomic.AddInt32(&i, 1)
		server.request.Async(remoteHost, remotePort, iMessage.InputMessage, func(replayMessage message.IMessage, hasReplay bool, err error) {
			atomic.AddInt32(&i, -1)
			if hasReplay{
				ty:=replayMessage.GetMessageType()
				if ty==message.BackMessageOKType{
					flag = true
					writeFunc(nil,true)
					server.userStore.AddUser(iMessage.GetToUsername(),remoteHost,remotePort)
				}
			}
			if !flag && i==0{
				writeFunc(err,false)
			}
		})
	})
	return hasMachine
}

func (server *Server) AcceptConn() {
	for {
		io, err := server.tcpserver.Accept()
		if err != nil {
			break
		} else {
			client, err := NewClient(io, server,server.context)
			if err == nil {
				go client.Start()
			}
		}
	}
}
func (server *Server) machineInfo(value ...interface{}) interface{} {
	mis := make([]*MachineInfo, 0)
	queryInfo:=newQueryMachineInfo()
	server.machineMap.eachAddress(func(remoteHost string, remotePort int) {
		im,_,err:= server.request.Call(remoteHost,remotePort,queryInfo)
		log.DebugF("查询machineInfo {}：{} err：{}",remoteHost,remotePort,err)
		if err==nil{
			data:=im.GetValue(message.QueryMachineInfo)
			if len(data)>0{
				var mi MachineInfo
				err1:=ffjson.Unmarshal(data,&mi)
				if err1==nil{
					mi.Address = remoteHost+":"+strconv.Itoa(remotePort)
					mis = append(mis,&mi)
				}else{
					log.InfoF("json 转换错误：{}",err1)
				}
			}
		}
	})
	vvv:=server.queryMachineInfo()
	mis = append(mis, (vvv).(*MachineInfo))
	return mis
}
func (server *Server) queryMachineInfo(value ...interface{})interface{}{
	var mi MachineInfo
	mi.Address = "localhost" + ":" + strconv.Itoa(server.port)
	mi.MachineId = server.machineId
	mi.UserNum = server.context.UserNum()
	return &mi
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
				log.ErrorF("queryMachineBasic err:{}", err)
			}
		}
		if hasQuery {
			server.getMachineList()
		}
		time.Sleep(time.Minute)
	}
}
func (server *Server) queryMachineBasic1(host string, port int) (*machine, *km.Conn, error) {
	qBasic := newQueryMachineBasic(server.port, server.machineId)
	msg, conn, err := server.request.Call(host, port, qBasic)
	if err == nil {
		machine := msg.GetString(message.BackMachineAddress)
		m, err := toMachine(machine)
		if err == nil {
			return m, conn, nil
		} else {
			return nil, nil, err
		}
	}
	return nil, nil, err
}

/**客户端获取请求机器信息**/
func (server *Server) queryMachineBasic() error {

	qBasic := newQueryMachineBasic(server.port, server.machineId)
	msg, conn, err := server.request.Call(server.machine.remoteHost, server.machine.remotePort, qBasic)
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
		log.DebugF("!!!!!!!getMachineList  remoteHost：{} remotePort：{}", remoteHost, remotePort)
		msg, _, err := server.request.Call(remoteHost, remotePort, qMsg)
		if err != nil {
			log.ErrorF("getMachineList err:{}  remoteHost：{} remotePort：{}", err, remoteHost, remotePort)
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
	/**
	集群用户记录
	 */

	server.userStore = newUserStore()
	server.port = context.GetConfig().GetIntOrDefault("cluster.local.port", 6361)
	server.machineId = context.GetConfig().GetStringOrDefault("cluster.local.machineId", MachineId())
	if server.machineId == "" {
		log.PanicF("machineId can‘t blank")
	}
	var remotePort = context.GetConfig().GetIntOrDefault("cluster.remote.port", 6362)
	var remoteHost = context.GetConfig().GetString("cluster.remote.host")
	if remoteHost == "" {
		log.PanicF("cluster.remote.host can‘t blank")
	}
	server.request = km.NewRequest()
	server.machineMap = newMachineStore()
	server.machine = newMachine(remotePort, remoteHost)
	server.tcpserver = net.NewTCPServer(server.port)
	context.HandleAddUser(server.HandleAddUser)
	context.HandleDeleteUser(server.HandleDeleteUser)
	context.HandleSendMessage(server.HandleSendMessage)
	context.HandleSendMultiMessage(server.HandleSendMultiMessage)
	context.QueryForwardHandle(server.Query)
	context.RegisterHandle("machineInfo", server.machineInfo)
	context.RegisterHandle("queryMachineInfo", server.queryMachineInfo)
	context.RegisterHandle("machineInfoId", server.machineInfoId)
}
func (server *Server) Name() string {
	return "cluster"
}

func (server *Server) machineInfoId(value ...interface{}) interface{} {
	return server.machineId
}

func NewServer() *Server {
	return &Server{}
}
