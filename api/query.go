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


func (query *Query) queryUser(value ...interface{}) interface{}{
	var u User
	machineInfoId:=query.context.GetHandle("machineInfoId")
	u.Machine=machineInfoId()
	u.RemoteAddress = make([]string,0)
	query.context.GetUser(value[0].(string), func(user user.IUser) bool {
		log.Info(user.GetUsername())
		u.Username = user.GetUsername()
		u.RemoteAddress = append(u.RemoteAddress, user.GetRemoteAddress())
		return true
	})
	return &u
}
func (query *Query) Init() {
	query.context.RegisterHandle("queryUser", query.queryUser)
	query.context.RegisterHandle("systemInfo", query.systemInfo)
}

type User struct {
	Username string
	RemoteAddress []string
	Machine interface{}
}

// NewValue go gob要求/**/
func (u *User)NewValue()interface{}  {
	var nu User
	nu.RemoteAddress = make([]string,0)
	return &nu
}