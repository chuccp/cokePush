package ex

import (
	log "github.com/chuccp/coke-log"
	"github.com/chuccp/cokePush/core"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/util"
	"net/http"
	"strconv"
	"sync"
	"time"
	"unsafe"
)

type store struct {
	clientMap *sync.Map
	context *core.Context
}

func (store *store)jack(w http.ResponseWriter, re *http.Request)  {
	userId:=util.GetUsername(re)
	log.DebugF("{} 来了",userId)
	v,ok:=store.clientMap.Load(userId)
	if ok{
		client:=v.(*client)
		client.poll(w)
	}else{
		store.createUser(userId,w)
	}
}
func (store *store)createUser(userId string,w http.ResponseWriter){
	client:=NewClient(store.context,userId)
	login:=message.CreateLoginMessage(userId)
	client.HandleLogin(login)
	store.clientMap.Store(userId,client)
	client.poll(w)
}

func (store *store)sendMsg(w http.ResponseWriter, re *http.Request)  {

}
func newStore(context *core.Context) *store {
	return &store{clientMap:new(sync.Map),context:context}
}

type client struct {
	queue *core.Queue
	context *core.Context
	username string
	userId string
}

func NewClient(context *core.Context,username string) *client {
	c:= &client{queue: core.NewQueue(),context:context,username:username}
	c.userId = username+strconv.FormatUint(uint64(uintptr(unsafe.Pointer(c))),36)
	return c
}
func (client *client)HandleLogin(iMessage message.IMessage){
	client.context.Handle(iMessage,client)
}
func (client *client)WriteMessage(iMessage message.IMessage) error{
	client.queue.Offer(iMessage)
	return nil
}
func (client *client)poll(w http.ResponseWriter) {
	msg:= client.queue.Poll(time.Second*20)
	if msg!=nil{
		w.Write(msg.GetValue(message.Text))
	}
}
func (client *client)GetUserId() string{
	return client.userId
}
func (client *client)GetUsername() string{
	return client.username
}
func (client *client)ReadMessage() (message.IMessage,error){
	return nil,nil
}
func (client *client)Close() {

}