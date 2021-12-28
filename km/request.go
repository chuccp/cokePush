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
	t:=time.Now()
	t = t.Add(time.Second*5)
	return &messageQ{writeMsg: msg, fa: fa.(chan bool),createTime: t,isWait:true,rLock:new(sync.RWMutex)}
}
func freeMessageQ(msg *messageQ) {
	poolChan.Put(msg.fa)
}

type messageQ struct {
	writeMsg message.IMessage
	readMsg  message.IMessage
	createTime time.Time
	isWait   bool
	fa       chan bool
	rLock *sync.RWMutex
}


func (q *messageQ) wait() (message.IMessage, bool) {
	return q.readMsg, <-q.fa
}

func (q *messageQ) notify(msg message.IMessage) {
	if !q.isWait{
		return}
	q.rLock.Lock()
	defer q.rLock.Unlock()
	if q.isWait{
		q.isWait = false
		q.readMsg = msg
		q.fa <- true
	}

}
func (q *messageQ) close() {
	if !q.isWait{
		return}
	q.rLock.Lock()
	defer q.rLock.Unlock()
	if q.isWait{
		q.isWait = false
		q.fa <- false
	}
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
		conn.msgChanMap.Store(msgId, msgQ)
		err:=conn.stream.WriteMessage(iMessage)
		if err!=nil{
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

func (conn *conn) getStatus() int {
	return conn.status
}
func (conn *conn) Close() {
	conn.stream.close(0)
	conn.msgChanMap.Range(func(key, value interface{}) bool {
		value.(*messageQ).close()
		return true
	})
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
			go conn.closeTimeOutMessage()
			go conn.read()
			go conn.live()
			conn.status = CONNING
		} else {
			conn.status = BREAK
		}
		return err
	}
}
/**
读取信息
 */
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
			}else{
				log.DebugF("读取到心跳反馈信息 classId:{}", classId)

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
/**
关掉超时消息
 */
func (conn *conn)closeTimeOutMessage()  {
	for conn.status == CONNING {
		time.Sleep(time.Second * 5)
		t:=time.Now()
		conn.msgChanMap.Range(func(key, value interface{}) bool {
			mq:=value.(*messageQ)
			if t.After(mq.createTime){
				mq.close()
			}
			return true
		})
	}
	
}

func newConn(host string, port int) *conn {
	return &conn{status: NEW, host: host, port: port, msgChanMap: new(sync.Map)}
}

type Request struct {
	connMap *sync.Map
	rLock   *util.MapLock
}

func NewRequest() *Request {
	return &Request{connMap: new(sync.Map), rLock: util.NewMapLock()}
}
func (request *Request) getConn(host string, port int) (*conn, error) {
	key := strconv.Itoa(port) + host
	val, ok := request.connMap.Load(key)
	if ok {
		conn := val.(*conn)
		return  request.connStatus(conn, key, host, port)
	} else {
		request.rLock.Lock(key)
		val, ok := request.connMap.Load(key)
		if ok {
			request.rLock.UnLock(key)
			conn := val.(*conn)
			return request.connStatus(conn, key, host, port)
		} else {
			nn, err := request.newConn(key, host, port)
			request.rLock.UnLock(key)
			return nn, err
		}
	}
	return nil, nil
}

func (request *Request) connStatus(conn *conn, key string, host string, port int) (*conn, error) {
	if conn.getStatus() == CONNING {
		return conn, nil
	}
	if conn.getStatus() == CREATING {
		return nil, core.ConnOnCreating
	}
	if conn.getStatus() == NEW ||conn.getStatus() == BREAK{
		request.rLock.Lock(key)
		if conn.getStatus() == CONNING {
			request.rLock.UnLock(key)
			return conn, nil
		}
		if conn.getStatus() == CREATING {
			request.rLock.UnLock(key)
			return nil, core.ConnOnCreating
		}
		if conn.getStatus() == NEW ||conn.getStatus() == BREAK{
			cnn,err:= request.newConn(key, host, port)
			request.rLock.UnLock(key)
			return cnn,err
		}
		request.rLock.UnLock(key)
	}
	request.rLock.UnLock(key)
	return nil, core.UnKnownConn
}

func (request *Request) newConn(key string, host string, port int) (*conn, error) {
	cn := newConn(host, port)
	request.connMap.Store(key, cn)
	err := cn.start()
	if err != nil {
		return nil, err
	}
	return cn, nil
}

func (request *Request) Call(host string, port int, message message.IMessage) (message.IMessage,*conn, error) {

	rq, err := request.getConn(host, port)
	if err != nil {
		return nil, rq,err
	} else {
		msg,err:=rq.write(message)
		return msg,rq,err
	}
}
