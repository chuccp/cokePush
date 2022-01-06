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
	rLock *sync.RWMutex
}

func (store *store) jack(w http.ResponseWriter, re *http.Request) {
	userId := util.GetUsername(re)
	v, ok := store.clientMap.Load(userId)
	if ok {
		ct := v.(*client)
		if !ct.poll(w){
			store.createUser(userId, w)
		}
	} else {
		store.createUser(userId, w)
	}
}
func (store *store) createUser(userId string, w http.ResponseWriter) {
	client := NewClient(store.context, userId)
	login := message.CreateLoginMessage(userId)
	client.HandleLogin(login)
	store.rLock.RLock()
	store.clientMap.Store(userId, client)
	store.rLock.RUnlock()
	client.poll(w)
}
func (store *store) deleteUser(userId string,c *client,t *time.Time) {
	flag:=c.close(t)
	if flag{
		store.rLock.Lock()
		v,ok:=store.clientMap.Load(userId)
		if ok{
			ci:=v.(*client)
			if ci==c{
				store.clientMap.Delete(userId)
			}
		}
		store.rLock.Unlock()
		store.context.DeleteUser(c)
	}
}
func (store *store) timeOutCheck() {
	for {
		time.Sleep(time.Second * 5)
		ti := time.Now()
		log.DebugF("扫描过期连接 {}", ti)
		store.clientMap.Range(func(key, value interface{}) bool {
			client := value.(*client)
			if client.timeOut(&ti) {
				store.deleteUser(key.(string),client,&ti)
			}
			return true
		})
	}
}
func (store *store) sendMsg(w http.ResponseWriter, re *http.Request) {

}
func newStore(context *core.Context) *store {
	return &store{clientMap: new(sync.Map), context: context,rLock:new(sync.RWMutex)}
}

type client struct {
	queue    *core.Queue
	context  *core.Context
	username string
	userId   string
	intPut   int
	hasClose bool
	last     *time.Time
	rLock *sync.RWMutex
}

func NewClient(context *core.Context, username string) *client {
	c := &client{queue: core.NewQueue(), context: context, username: username, intPut: 0,hasClose:false,rLock:new(sync.RWMutex)}
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
func (client *client)close(t *time.Time)bool{
	client.rLock.Lock()
	if client.timeOut(t){
		client.hasClose = true
		client.rLock.Unlock()
		return true
	}
	client.rLock.Unlock()
	return false
}

func (client *client)isClose() bool {
	return client.hasClose
}

func (client *client) poll(w http.ResponseWriter) bool {
	client.rLock.RLock()
	if client.hasClose{
		client.rLock.RUnlock()
		return false
	}
	client.intPut++
	client.rLock.RUnlock()
	msg := client.queue.Poll(time.Second * 20)
	if msg != nil {
		w.Write(msg.GetValue(message.Text))
	}
	client.rLock.RLock()
	client.intPut--
	if client.intPut == 0 {
		t := time.Now().Add(time.Second * 10)
		client.last = &t
	}
	client.rLock.RUnlock()
	return true
}
func (client *client) timeOut(t *time.Time) bool {
	if client.intPut == 0 {
		flag := client.last.Before(*t)
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
