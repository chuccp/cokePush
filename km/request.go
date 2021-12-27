package km

import (
	log "github.com/chuccp/coke-log"
	"github.com/chuccp/cokePush/core"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/net"
	"strconv"
	"sync"
	"time"
)

const (
	NEW = iota
	CREATING
	CONNING
	BREAK
)

var poolChan = &sync.Pool{
	New: func() interface{} {
		return make(chan bool)
	},
}

func getMessageQ(msg message.IMessage) *messageQ {
	fa := poolChan.Get()
	return &messageQ{writeMsg: msg, fa: fa.(chan bool)}
}
func freeMessageQ(msg *messageQ) {
	poolChan.Put(msg.fa)
}

type messageQ struct {
	writeMsg message.IMessage
	readMsg  message.IMessage
	fa       chan bool
}

func newMessageQ(msg message.IMessage) *messageQ {
	return &messageQ{writeMsg: msg, fa: make(chan bool)}
}
func (q *messageQ) wait() (message.IMessage, bool) {
	return q.readMsg, <-q.fa
}

func (q *messageQ) notify(msg message.IMessage) {
	q.readMsg = msg
	q.fa <- true
}
func (q *messageQ) close() {
	q.fa <- false
}

type conn struct {
	status     int
	host       string
	port       int
	stream     *Stream
	msgChanMap *sync.Map
}

func (conn *conn) write(iMessage message.IMessage) (message.IMessage, error) {
	msgId := iMessage.GetMessageId()
	if msgId == 0 {
		return nil, core.MessageIdIsBlank
	} else {
		msgQ := getMessageQ(iMessage)
		conn.msgChanMap.Store(msgId, iMessage)
		msg, fa := msgQ.wait()
		freeMessageQ(msgQ)
		if fa {
			return msg, nil
		} else {
			return nil, core.ReadTimeout
		}
	}
	return nil, nil

}

func (conn *conn) getStatus() int {
	return conn.status
}

func (conn *conn) start() error {
	conn.status = CREATING
	c := net.NewXConn(conn.host, conn.port)
	sm, err := c.Create()
	if err != nil {
		conn.status = BREAK
		return err
	} else {
		conn.stream, err = NewClientStream(sm)
		if err == nil {
			go conn.read()
			go conn.live()
			conn.status = CONNING
		} else {
			conn.status = BREAK
		}
		return err
	}
}
func (conn *conn) read() {
	for {
		msg, err := conn.stream.ReadMessage()
		if err == nil {
			classId := msg.GetClassId()
			if classId != message.LiveMessageClass {
				msgId := msg.GetMessageId()
				msgQ, ok := conn.msgChanMap.Load(msgId)
				if ok {
					mq, ok := msgQ.(*messageQ)
					if ok {
						mq.notify(msg)
					}
				} else {
					log.InfoF("读取到超时反馈信息 classId:{} msgId:{}", classId, msgId)
				}
			}
		} else {
			break
		}
	}
	conn.status = BREAK
}

/**心跳维持连接**/
func (conn *conn) live() {

	for conn.status == CONNING {
		lm := message.CreateLiveMessage()
		err := conn.stream.WriteMessage(lm)
		if err != nil {
			break
		}
		time.Sleep(time.Minute * 10)
	}
	conn.status = BREAK
}

func newConn(host string, port int) *conn {
	return &conn{status: NEW, host: host, port: port, msgChanMap: new(sync.Map)}
}

type Request struct {
	connMap *sync.Map
	rLock   *sync.RWMutex
}

func NewRequest() *Request {
	return &Request{connMap: new(sync.Map), rLock: new(sync.RWMutex)}
}
func (request *Request) getConn(host string, port int) (*conn, error) {
	key := strconv.Itoa(port) + host
	request.rLock.RLock()
	val, ok := request.connMap.Load(key)
	if ok {
		conn := val.(*conn)
		if conn.getStatus() == CONNING {
			return conn, nil
		} else if conn.getStatus() == NEW {
			cn := newConn(host, port)
			return cn, nil
		} else if conn.getStatus() == BREAK {
			cn := newConn(host, port)
			return cn, nil
		} else if conn.getStatus() == CREATING {
			return nil, nil
		}
	} else {
		cn := newConn(host, port)
		request.connMap.Store(key, cn)
		return cn, nil
	}
	return nil, nil
}
func (request *Request) Call(host string, port int, message message.IMessage) (*message.IMessage, error) {

	return nil, nil
}
