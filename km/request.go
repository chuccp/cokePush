package km

import (
	log "github.com/chuccp/coke-log"
	"github.com/chuccp/cokePush/core"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/net"
	"github.com/chuccp/cokePush/util"
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
	t := time.Now()
	t = t.Add(time.Second * 5)
	return &messageQ{writeMsg: msg, fa: fa.(chan bool), createTime: t, isWait: true, rLock: new(sync.RWMutex)}
}
func freeMessageQ(msg *messageQ) {
	poolChan.Put(msg.fa)
}

type messageQ struct {
	writeMsg   message.IMessage
	readMsg    message.IMessage
	createTime time.Time
	isWait     bool
	fa         chan bool
	rLock      *sync.RWMutex
}
type messageB struct {
	writeMsg     message.IMessage
	createTime   time.Time
	callBackFunc CallBackFunc
}

func getMessageB(msg message.IMessage, callBackFunc CallBackFunc) *messageB {
	t := time.Now()
	t = t.Add(time.Second * 5)
	return &messageB{writeMsg: msg, createTime: t, callBackFunc: callBackFunc}
}

func (q *messageQ) wait() (message.IMessage, bool) {
	return q.readMsg, <-q.fa
}

func (q *messageQ) notify(msg message.IMessage) {
	if !q.isWait {
		return
	}
	q.rLock.Lock()
	defer q.rLock.Unlock()
	if q.isWait {
		q.isWait = false
		q.readMsg = msg
		q.fa <- true
	}

}
func (q *messageQ) close() {
	if !q.isWait {
		return
	}
	q.rLock.Lock()
	defer q.rLock.Unlock()
	if q.isWait {
		q.isWait = false
		q.fa <- false
	}
}

type Conn struct {
	status     int
	host       string
	port       int
	stream     *Stream
	msgChanMap *sync.Map
}

func (conn *Conn) write(iMessage message.IMessage) (message.IMessage, error) {
	msgId := iMessage.GetMessageId()
	if msgId == 0 {
		return nil, core.MessageIdIsBlank
	} else {
		msgQ := getMessageQ(iMessage)
		conn.msgChanMap.Store(msgId, msgQ)
		err := conn.stream.WriteMessage(iMessage)
		if err != nil {
			conn.msgChanMap.Delete(msgId)
			freeMessageQ(msgQ)
			return nil, err
		}
		log.DebugF("==========开始等待")
		msg, fa := msgQ.wait()
		log.DebugF("==========等待结束")
		conn.msgChanMap.Delete(msgId)
		freeMessageQ(msgQ)
		if fa {
			return msg, nil
		} else {
			return nil, core.ReadTimeout
		}
	}
	return nil, nil

}

type CallBackFunc func(message.IMessage, bool, error)

func (conn *Conn) asyncWrite(iMessage message.IMessage, callBackFunc CallBackFunc) {
	msgId := iMessage.GetMessageId()
	if msgId == 0 {
		callBackFunc(nil, false, core.MessageIdIsBlank)
	} else {
		msgB := getMessageB(iMessage, callBackFunc)
		conn.msgChanMap.Store(msgId, msgB)
		err := conn.stream.WriteMessage(iMessage)
		if err != nil {
			conn.msgChanMap.Delete(msgId)
			callBackFunc(nil, false, err)
		}
	}
}

func (conn *Conn) getStatus() int {
	return conn.status
}
func (conn *Conn) Close() {
	conn.stream.close(0)
}
func (conn *Conn) clear() {
	conn.msgChanMap.Range(func(key, value interface{}) bool {


		switch mq := value.(type) {
		case *messageQ:
			{
				mq.close()
			}
		case *messageB:
			{
				conn.msgChanMap.Delete(mq.writeMsg.GetMessageId())
			}
		}
		return true
	})
}
func (conn *Conn) start() error {
	conn.status = CREATING
	c := net.NewXConn(conn.host, conn.port)
	sm, err := c.Create()
	if err != nil {
		conn.status = BREAK
		return err
	} else {
		conn.stream, err = NewClientStream(sm)
		if err == nil {
			conn.status = CONNING
			go conn.closeTimeOutMessage()
			go conn.read()
			go conn.live()
		} else {
			conn.status = BREAK
			conn.stream.close(0)
		}
		return err
	}
}

/**
读取信息
*/
func (conn *Conn) read() {
	for conn.status == CONNING {
		msg, err := conn.stream.ReadMessage()
		if err == nil {
			log.InfoF("收到信息：classId:{} type:{} msgId:{}", msg.GetClassId(), msg.GetMessageType(), msg.GetMessageId())
			classId := msg.GetClassId()
			if classId != message.LiveMessageClass {
				msgId := msg.GetMessageId()
				ms, ok := conn.msgChanMap.Load(msgId)
				if ok {
					switch mq := ms.(type) {
					case *messageQ:
						{
							log.InfoF("messageQ msgId:{}", msgId)
							mq.notify(msg)
						}
					case *messageB:
						{
							log.InfoF("messageB msgId:{}", msgId)
							conn.msgChanMap.Delete(msgId)
							mq.callBackFunc(msg, true, nil)
						}
					}
				} else {
					log.InfoF("读取到超时反馈信息 classId:{} msgId:{}", classId, msgId)
				}
			} else {
				log.DebugF("读取到心跳反馈信息 classId:{}", classId)

			}
		} else {
			break
		}
	}
	conn.status = BREAK
	conn.clear()
}

/**心跳维持连接**/
func (conn *Conn) live() {

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

/**
关掉超时消息
*/
func (conn *Conn) closeTimeOutMessage() {
	for conn.status == CONNING {
		log.DebugF("扫描过期消息=======")
		time.Sleep(time.Second * 5)
		t := time.Now()
		conn.msgChanMap.Range(func(key, value interface{}) bool {
			switch mq := value.(type) {
			case *messageQ:
				{
					if t.After(mq.createTime) {
						mq.close()
					}
				}
			case *messageB:
				{
					conn.msgChanMap.Delete(mq.writeMsg.GetMessageId())
					mq.callBackFunc(nil, false, TimeoutError)
				}
			}
			return true
		})
	}
	log.DebugF("扫描过期消息=======结束")
}

func newConn(host string, port int) *Conn {
	return &Conn{status: NEW, host: host, port: port, msgChanMap: new(sync.Map)}
}

type Request struct {
	connMap *sync.Map
	rLock   *util.MapLock
}

func NewRequest() *Request {
	return &Request{connMap: new(sync.Map), rLock: util.NewMapLock()}
}
func (request *Request) getConn(host string, port int) (*Conn, error) {
	key := strconv.Itoa(port) + host
	val, ok := request.connMap.Load(key)
	if ok {
		conn := val.(*Conn)
		return request.connStatus(conn, key, host, port)
	} else {
		request.rLock.Lock(key)
		val, ok := request.connMap.Load(key)
		if ok {
			request.rLock.UnLock(key)
			conn := val.(*Conn)
			return request.connStatus(conn, key, host, port)
		} else {
			nn, err := request.newConn(key, host, port)
			request.rLock.UnLock(key)
			return nn, err
		}
	}
	return nil, nil
}

func (request *Request) connStatus(conn *Conn, key string, host string, port int) (*Conn, error) {
	if conn.getStatus() == CONNING {
		return conn, nil
	}
	if conn.getStatus() == CREATING {
		return nil, core.ConnOnCreating
	}
	if conn.getStatus() == NEW || conn.getStatus() == BREAK {
		request.rLock.Lock(key)
		if conn.getStatus() == CONNING {
			request.rLock.UnLock(key)
			return conn, nil
		}
		if conn.getStatus() == CREATING {
			request.rLock.UnLock(key)
			return nil, core.ConnOnCreating
		}
		if conn.getStatus() == NEW || conn.getStatus() == BREAK {
			cnn, err := request.newConn(key, host, port)
			request.rLock.UnLock(key)
			return cnn, err
		}
		request.rLock.UnLock(key)
	}
	request.rLock.UnLock(key)
	return nil, core.UnKnownConn
}

func (request *Request) newConn(key string, host string, port int) (*Conn, error) {
	cn := newConn(host, port)
	request.connMap.Store(key, cn)
	err := cn.start()
	if err != nil {
		return nil, err
	}
	return cn, nil
}

func (request *Request) Call(host string, port int, message message.IMessage) (message.IMessage, *Conn, error) {

	rq, err := request.getConn(host, port)
	if err != nil {
		return nil, rq, err
	} else {
		msg, err := rq.write(message)
		return msg, rq, err
	}
}
func (request *Request) Async(host string, port int, iMessage message.IMessage, callBackFunc CallBackFunc) {
	rq, err := request.getConn(host, port)
	if err != nil {
		callBackFunc(nil, false, err)
	} else {
		rq.asyncWrite(iMessage, func(iMessage message.IMessage, b bool, err error) {
			callBackFunc(iMessage, b, err)
		})
	}
}
