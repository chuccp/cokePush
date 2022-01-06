package ex

import (
	log "github.com/chuccp/coke-log"
	"github.com/chuccp/cokePush/core"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/util"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
	"unsafe"
)

type store struct {
	clientMap *sync.Map
	context   *core.Context
}

func (store *store) jack(w http.ResponseWriter, re *http.Request) {
	userId := util.GetUsername(re)
	v, ok := store.clientMap.Load(userId)
	if ok {
		client := v.(*client)
		client.poll(w)
	} else {
		store.createUser(userId, w)
	}
}
func (store *store) createUser(userId string, w http.ResponseWriter) {
	client := NewClient(store.context, userId)
	login := message.CreateLoginMessage(userId)
	client.HandleLogin(login)
	store.clientMap.Store(userId, client)
	client.poll(w)
}
func (store *store) deleteUser(userId string,client *client) {
	client.hasClose = true
	store.clientMap.Delete(userId)
	store.context.DeleteUser(client)
}
func (store *store) timeOutCheck() {
	for {
		time.Sleep(time.Second * 5)
		ti := time.Now()
		log.DebugF("扫描过期连接 {}", ti)
		store.clientMap.Range(func(key, value interface{}) bool {
			client := value.(*client)
			if client.timeOut(&ti) {
				store.deleteUser(key.(string),client)
			}
			return true
		})
	}
}
func (store *store) sendMsg(w http.ResponseWriter, re *http.Request) {

}
func newStore(context *core.Context) *store {
	return &store{clientMap: new(sync.Map), context: context}
}

type client struct {
	queue    *core.Queue
	context  *core.Context
	username string
	userId   string
	intPut   int
	hasClose bool
	last     *time.Time
}

func NewClient(context *core.Context, username string) *client {
	c := &client{queue: core.NewQueue(), context: context, username: username, intPut: 0,hasClose:false}
	c.userId = username + strconv.FormatUint(uint64(uintptr(unsafe.Pointer(c))), 36)
	return c
}
func (client *client) HandleLogin(iMessage message.IMessage) {
	client.context.Handle(iMessage, client)
}
func (client *client) WriteMessage(iMessage message.IMessage) error {
	if client.hasClose{
		return net.ErrClosed
	}
	log.DebugF("WriteMessage messageId {}", iMessage.GetMessageId())
	client.queue.Offer(iMessage)
	return nil
}
func (client *client) poll(w http.ResponseWriter) {
	client.intPut++
	msg := client.queue.Poll(time.Second * 20)
	if msg != nil {
		w.Write(msg.GetValue(message.Text))
	}
	client.intPut--
	if client.intPut == 0 {
		t := time.Now().Add(time.Second * 10)
		client.last = &t
	}
}
func (client *client) timeOut(t *time.Time) bool {
	if client.intPut == 0 {
		flag := client.last.Before(*t)
		log.DebugF("{} 超时检查用户  timeOut:{} nowTime:{} 是否超时：{}", client.username, client.last, t, flag)
		return flag
	}
	return false
}
func (client *client) getLastTime() *time.Time {
	if client.intPut != 0 {
		return nil
	}
	return client.last
}

func (client *client) GetId() string {
	return client.userId
}
func (client *client) GetUsername() string {
	return client.username
}
func (client *client)SetUsername(username string){
	client.username = username
}
func (client *client) ReadMessage() (message.IMessage, error) {
	return nil, nil
}
