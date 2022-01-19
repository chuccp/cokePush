package ex

import (
	"context"
	log "github.com/chuccp/coke-log"
	"github.com/chuccp/cokePush/core"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/user"
	"github.com/chuccp/cokePush/util"
	"github.com/chuccp/queue"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type client struct {
	queue    *queue.Queue
	context  *core.Context
	username string
	rLock   *sync.RWMutex
	connMap *sync.Map
	connNum int32
	last    *time.Time
	intPut   int32
}
type conn struct {
	userId     string
	username   string
	address    string
	w          http.ResponseWriter
	re         *http.Request
	client     *client
	add        *time.Time
	last       *time.Time
	isWrite    int32
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func newConn(username string, w http.ResponseWriter, re *http.Request, client *client) *conn {
	c := &conn{w: w, re: re}
	c.SetUsername(username)
	c.address = re.RemoteAddr
	c.userId = username + re.RemoteAddr
	c.client = client
	return c
}
func (u *conn) WriteMessage(iMessage message.IMessage) error {
	u.client.queue.Offer(iMessage.GetValue(message.Text))
	return nil
}
func (u *conn) GetId() string {
	return u.userId
}
func (u *conn) GetUsername() string {
	return u.username
}
func (u *conn) GetRemoteAddress() string {
	return u.address
}
func (u *conn) SetUsername(username string) {
	u.username = username
}
func (u *conn) writeBlank() {
	if u.canWrite() {
		u.w.Write([]byte("[]"))
		u.cancelFunc()
	}
}
func (u *conn) WriteMessageFunc(iMessage message.IMessage, writeFunc user.WriteFunc) {

}
func (u *conn) toWrite() {
	atomic.StoreInt32(&u.isWrite, 1)
}
func (u *conn) canWrite() bool {
	return atomic.CompareAndSwapInt32(&u.isWrite, 1, -1)
}

func (c *client) poll(username string, w http.ResponseWriter, re *http.Request) {
	atomic.AddInt32(&c.intPut, 1)
	cnn := newConn(username, w, re, c)
	v, flag := c.connMap.LoadOrStore(cnn.userId, cnn)
	cnn = v.(*conn)
	cnn.toWrite()
	if !flag {
		c.connNum++
		c.context.AddUser(cnn)
	}
	ti := time.Now().Add(time.Second * 25)
	cnn.add = &ti
	cnn.ctx, cnn.cancelFunc = context.WithTimeout(context.Background(), time.Minute)
	v, _, cls := c.queue.Dequeue(cnn.ctx)
	if cnn.canWrite() {
		if cls {
			cnn.w.Write([]byte("[]"))
		} else {
			cnn.cancelFunc()
			if v != nil {
				cnn.w.Write(v.([]byte))
			}
		}
	} else {
		if v != nil {
			c.queue.Offer(v)
		}
	}
	t := time.Now().Add(time.Second * 10)
	num := atomic.AddInt32(&c.intPut, -1)
	if num == 0 {
		c.last = &t
	}
	cnn.last = &t
}
func (c *client) timeoutCheck(t *time.Time) {
	c.connMap.Range(func(key, value interface{}) bool {
		cnn := value.(*conn)
		if cnn.last != nil && cnn.last.Before(*t) {
			log.InfoF("超时 {}   {}", cnn.username, cnn.userId)
			c.context.DeleteUser(cnn)
			c.connMap.Delete(cnn.userId)
			c.connNum--
		}
		return true
	})
}

func (c *client) writeBlank(t *time.Time) {
	c.connMap.Range(func(key, value interface{}) bool {
		cnn := value.(*conn)
		if cnn.add != nil && cnn.add.Before(*t) {
			cnn.writeBlank()
		}
		return true
	})
}
func newClient(context *core.Context, username string) *client {
	c := &client{queue: queue.NewQueue(), context: context, username: username, rLock: new(sync.RWMutex), connMap: new(sync.Map)}
	return c
}

type store struct {
	clientMap *sync.Map
	context   *core.Context
	rLock     *sync.RWMutex
}

func (store *store) jack(w http.ResponseWriter, re *http.Request) {
	username := util.GetUsername(re)
	cl := newClient(store.context, username)
	v, _ := store.clientMap.LoadOrStore(username, cl)
	ct := v.(*client)
	ct.poll(username, w, re)
}
func (store *store) timeoutCheck() {
	for {
		time.Sleep(time.Second * 10)
		ti := time.Now()
		log.TraceF("扫描过期连接 {}", ti)
		store.clientMap.Range(func(key, value interface{}) bool {
			cl := value.(*client)
			cl.timeoutCheck(&ti)
			if cl.connNum==0{
				store.clientMap.Delete(key)
			}
			return true
		})
	}
}

func (store *store) writeBlank() {
	log.InfoF("轮询检查http长链接")
	for {
		time.Sleep(time.Second * 10)
		t := time.Now()
		store.clientMap.Range(func(key, value interface{}) bool {
			c, ok := value.(*client)
			if ok {
				c.writeBlank(&t)
			}
			return true
		})
	}
}

func (store *store) sendMsg(w http.ResponseWriter, re *http.Request) {

}
func newStore(context *core.Context) *store {
	return &store{clientMap: new(sync.Map), context: context, rLock: new(sync.RWMutex)}
}
