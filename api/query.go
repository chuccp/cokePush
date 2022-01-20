package api

import (
	log "github.com/chuccp/coke-log"
	"github.com/chuccp/cokePush/core"
	"github.com/chuccp/cokePush/user"
)

const (
	VERSION = "version"
)

type Query struct {
	context *core.Context
}

func newQuery(context *core.Context) *Query {
	q:= &Query{context:context}
	q.Init()
	return q
}

type SystemInfo struct {
	SendMsgNum int32
	ReplayMsgNum int32
	Machine interface{}
}
// NewValue go gob要求/**/
func (u *SystemInfo)NewValue()interface{}  {
	var nu SystemInfo
	return &nu
}
func (query *Query) systemInfo(value ...interface{})interface{}{
	var systemInfo SystemInfo
	systemInfo.SendMsgNum = query.context.SendNum()
	systemInfo.ReplayMsgNum = query.context.ReplyNum()
	machineInfoId:=query.context.GetHandle("machineInfoId")
	systemInfo.Machine = machineInfoId()
	return &systemInfo
}

func (query *Query) queryPageUser(value ...interface{}) interface{} {
	start:=value[0].(int)
	size:=value[1].(int)

	log.InfoF("queryPageUser  start:{} size:{}",start,size)

	page:=user.NewPage()
	page.Num = (int)(query.context.UserNum())
	if start>=page.Num{
		return page
	}
	var num int = 0
	query.context.EachUsers(func(key string, value *user.StoreUser) bool {

		if num==start{
			if size==0{
				return false
			}
			size--
			page.List = append(page.List, user.NewPageUser(value.GetUsername(),value.MachineAddress(),value.CreateTime()))

		}
		num++
		return true
	})
	return page
}

func (query *Query) queryUser(value ...interface{}) interface{}{
	var u User
	machineInfoId:=query.context.GetHandle("machineInfoId")
	u.Machine=machineInfoId()
	u.Conn = make([]*Conn,0)
	query.context.GetUser(value[0].(string), func(user user.IUser) bool {
		log.Info(user.GetUsername())
		u.Username = user.GetUsername()
		u.Conn = append(u.Conn, newConn(user.GetRemoteAddress(),"",user.GetRemoteAddress()))
		return true
	})
	return &u
}
func (query *Query) Init() {
	query.context.RegisterHandle("queryUser", query.queryUser)
	query.context.RegisterHandle("systemInfo", query.systemInfo)
	query.context.RegisterHandle("queryPageUser", query.queryPageUser)
}

type User struct {
	Username string
	Conn []*Conn
	Machine interface{}
}
type Conn struct {
	RemoteAddress string
	LastLiveTime string
	CreateTime string
}

func newConn(RemoteAddress string,LastLiveTime string,CreateTime string)*Conn{
	return &Conn{RemoteAddress,LastLiveTime,CreateTime}
}