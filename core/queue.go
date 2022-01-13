package core

import (
	"container/list"
	"github.com/chuccp/cokePush/message"
	"github.com/chuccp/cokePush/util"
	"sync"
	"time"
)

const (
	wait int8 = iota
	run
)

type DockMessage struct {
	InputMessage message.IMessage
	write        WriteFunc
	flag         bool
	err          error
	IsForward    bool
}

func (dm *DockMessage) GetToUsername() string {
	return dm.InputMessage.GetString(message.ToUser)
}
func newDockMessage(inputMessage message.IMessage, write WriteFunc) *DockMessage {
	return &DockMessage{InputMessage: inputMessage, flag: false, write: write}
}

type Queue struct {
	messageList *list.List
	ch          chan bool
	status      int8
	timer       *util.Timer
	lock        *sync.RWMutex
}

func NewQueue() *Queue {
	return &Queue{messageList: list.New(), ch: make(chan bool), lock: new(sync.RWMutex)}
}
func (queue *Queue) Num() int {
	return queue.messageList.Len()
}

func (queue *Queue) offer(msg *DockMessage) {
	queue.lock.Lock()
	queue.messageList.PushBack(msg)
	if queue.status == wait {
		queue.ch <- true
		queue.status = run
		queue.lock.Unlock()
	} else {
		queue.lock.Unlock()
	}
}
func (queue *Queue) Offer(msg message.IMessage) {
	queue.lock.Lock()
	queue.messageList.PushBack(msg)
	if queue.status == wait {
		queue.ch <- true
		queue.status = run
		queue.lock.Unlock()
	} else {
		queue.lock.Unlock()
	}
}
func (queue *Queue) poll() *DockMessage {
	for {
		queue.lock.Lock()
		ele := queue.messageList.Front()
		if ele != nil {
			queue.messageList.Remove(ele)
			queue.lock.Unlock()
			return ele.Value.(*DockMessage)
		} else {
			queue.status = wait
			queue.lock.Unlock()
			<-queue.ch
		}
	}
	return nil
}
func (queue *Queue) Poll(duration time.Duration) message.IMessage {
	for {
		queue.lock.Lock()
		ele := queue.messageList.Front()
		if ele != nil {
			queue.messageList.Remove(ele)
			queue.lock.Unlock()
			return ele.Value.(message.IMessage)
		} else {
			queue.status = wait
			queue.lock.Unlock()
			queue.timer = util.GetTimer(duration)
			go func() {
				fa := queue.timer.Wait()
				if !fa {
					queue.Offer(message.CreateBasicMessage("", "", "[]"))
				}
			}()
			<-queue.ch
			queue.timer.End()
			util.FreeTimer(queue.timer)
		}
	}
	return nil
}
