package km

import (
	log "github.com/chuccp/coke-log"
	"github.com/chuccp/cokePush/core"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/net"
	net2 "net"
	"strconv"
	"sync"
	"time"
)

type STATUS int

const (
	NEW STATUS = iota
	CREATING
	CONNING
	BREAK
	ERROR
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
	status     STATUS
	address       string
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
		msg, fa := msgQ.wait()
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

// JustWrite /** 只写不管有无处理
func (conn *Conn) justWrite(iMessage message.IMessage) {
	conn.stream.WriteMessage(iMessage)
}
func (conn *Conn) getStatus() STATUS {
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
				if mq.callBackFunc != nil {
					mq.callBackFunc(nil, false, net2.ErrClosed)
				}
				conn.msgChanMap.Delete(mq.writeMsg.GetMessageId())
			}
		}
		return true
	})
}
func (conn *Conn) start() error {
	conn.status = CREATING
	c := net.NewXConn2(conn.address)
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
			classId := msg.GetClassId()
			if classId != message.LiveMessageClass {
				msgId := msg.GetMessageId()
				ms, ok := conn.msgChanMap.Load(msgId)
				if ok {
					switch mq := ms.(type) {
					case *messageQ:
						{
							log.DebugF("messageQ msgId:{}", msgId)
							mq.notify(msg)
						}
					case *messageB:
						{
							log.DebugF("messageB msgId:{}", msgId)
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


func newConn(address string) *Conn {
	return &Conn{status: NEW, address: address, msgChanMap: new(sync.Map)}
}
type Request struct {
	connMap *sync.Map
	rLock   *sync.RWMutex
}

func NewRequest() *Request {
	return &Request{connMap: new(sync.Map), rLock: new(sync.RWMutex)}
}
func (request *Request) async(address string, f func(*Conn, STATUS, error)){
	val, ok := request.connMap.Load(address)
	if ok {

		conn := val.(*Conn)
		request.rLock.Lock()
		if conn.status == NEW || conn.status == BREAK {
			conn.status = CREATING
			request.rLock.Unlock()
			connStart(conn, f)
		} else {
			request.rLock.Unlock()
			f(conn, conn.status, nil)
		}
	} else {
		request.rLock.Lock()
		val, ok = request.connMap.Load(address)
		if ok {
			conn := val.(*Conn)
			if conn.status == NEW || conn.status == BREAK {
				conn.status = CREATING
				request.rLock.Unlock()
				connStart(conn, f)
			} else {
				request.rLock.Unlock()
				f(conn, conn.status, nil)
			}
		} else {
			cn := newConn(address)
			request.connMap.Store(address, cn)
			cn.status = CREATING
			request.rLock.Unlock()
			connStart(cn, f)
		}
	}
}
func (request *Request) async2(host string, port int, f func(*Conn, STATUS, error)) {
	request.async(host+":"+strconv.Itoa(port),f)
}

func connStart(cn *Conn, f func(*Conn, STATUS, error)) {
	go func() {
		err := cn.start()
		if err == nil {
			f(cn, cn.status, nil)
		} else {
			f(cn, ERROR, err)
		}
	}()
}
func (request *Request) Call(host string, port int, msg message.IMessage) (iMsg message.IMessage, cnn *Conn, err error) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	request.async2(host, port, func(conn *Conn, status STATUS, err1 error) {
		cnn = conn
		if status == CONNING {
			conn.asyncWrite(msg, func(iMessage message.IMessage, b bool, err2 error) {
				iMsg = iMessage
				err = err2
				if !b {
					err = core.ReadTimeout
				}
				wg.Done()
			})
		} else {
			err = err1
			wg.Done()
		}
	})
	wg.Wait()
	return
}
func (request *Request) Async(host string, port int, iMessage message.IMessage, callBackFunc CallBackFunc) {

	request.async2(host, port, func(conn *Conn, status STATUS, err error) {
		if status == CONNING {
			conn.asyncWrite(iMessage,callBackFunc)
		} else if status==CREATING{
			callBackFunc(nil, false,core.ConnOnCreating)
		}else{
			callBackFunc(nil, false,err)
		}
	})
}
func (request *Request) JustCall2(remoteAddress string, message message.IMessage) {
	request.async(remoteAddress, func(conn *Conn, status STATUS, err error) {
		if status == CONNING {
			conn.justWrite(message)
		}
	})
}
func (request *Request) JustCall(host string, port int, message message.IMessage) {
	request.async2(host, port, func(conn *Conn, status STATUS, err error) {
		if status == CONNING {
			conn.justWrite(message)
		}
	})
}
