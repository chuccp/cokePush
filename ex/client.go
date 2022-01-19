package ex

import (
	"github.com/chuccp/cokePush/core"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/user"
	"github.com/chuccp/cokePush/util"
	"github.com/chuccp/queue"
	"net/http"
	"sync"
	"time"
)

type client struct {
	queue    *queue.Queue
	context  *core.Context
	username string
	rLock    *sync.RWMutex
	userMap  *sync.Map
}
type conn struct {
	userId   string
	username string
	address  string
	w        http.ResponseWriter
	re       *http.Request
	client   *client
	add      *time.Time
	isWrite  bool
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
	u.client.queue.Offer(iMessage.GetString(message.Text))
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
	if !u.isWrite {
		u.isWrite = true
		u.w.Write([]byte("[]"))
		u.client.queue.Offer(nil)
	}
}
func (u *conn) WriteMessageFunc(iMessage message.IMessage, writeFunc user.WriteFunc) {

}
func (c *client) poll(username string, w http.ResponseWriter, re *http.Request) bool {

	cnn := newConn(username, w, re, c)
	cnn.isWrite = false
	_, flag := c.userMap.LoadOrStore(cnn.userId, cnn)
	if !flag {
		c.context.AddUser(cnn)
	}
	ti := time.Now()
	cnn.add = &ti
	v, _ := c.queue.Take(time.Second * 40)
	if !cnn.isWrite {
		cnn.isWrite = true
		if v != nil {
			m := v.(message.IMessage)
			cnn.w.Write(m.GetValue(message.Text))
		} else {
			cnn.w.Write([]byte("[]"))
		}
	}
	c.userMap.Delete(cnn.userId)
	return false
}

func newClient(context *core.Context, username string) *client {
	c := &client{queue: queue.NewQueue(), context: context, username: username, rLock: new(sync.RWMutex), userMap: new(sync.Map)}
	return c
}

type store struct {
	clientMap *sync.Map
	context   *core.Context
	rLock     *sync.RWMutex
}

func (store *store) jack(w http.ResponseWriter, re *http.Request) {
	username := util.GetUsername(re)
	v, ok := store.clientMap.Load(username)
	if ok {
		ct := v.(*client)
		if !ct.poll(username, w, re) {

		}
	} else {

	}
}
func (store *store) createUser(userId string, w http.ResponseWriter, re *http.Request) {

}
func (store *store) deleteUser(userId string, c *client, t *time.Time) {

}
func (store *store) timeOutCheck() {

}

func (store *store) writeBlank() {

}

func (store *store) sendMsg(w http.ResponseWriter, re *http.Request) {

}
func newStore(context *core.Context) *store {
	return &store{clientMap: new(sync.Map), context: context, rLock: new(sync.RWMutex)}
}

type client2 struct {
}
