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
func (query *Query) queryUser(value ...interface{}) interface{} {
	var u User
	machineInfoId:=query.context.GetHandle("machineInfoId")
	u.Machine=machineInfoId()
	query.context.GetUser(value[0].(string), func(user user.IUser) bool {
		log.Info(user.GetUsername())
		u.Username = user.GetUsername()
		u.Id = user.GetId()
		return true
	})
	return u
}
func (query *Query) Init() {
	query.context.RegisterHandle("QueryUser", query.queryUser)
}

type User struct {
	Username string
	Id string
	Machine interface{}
}
