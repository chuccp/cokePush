package core

import (
	"container/list"
	"github.com/chuccp/cokePush/message"
	"sync"
)

const (

	wait int8 = iota
	run
)

type queue struct {
	messageList *list.List
	ch chan bool
	status int8
	lock *sync.RWMutex
}
func newQueue()*queue{
	return &queue{messageList: list.New(),ch:make(chan bool),lock:new(sync.RWMutex)}
}
func (queue* queue)Num()int{
	return queue.messageList.Len()
}

func (queue* queue) offer(msg message.IMessage)  {
	queue.lock.RLock()
	queue.messageList.PushBack(msg)
	if queue.status==wait{
		queue.lock.RUnlock()
		queue.lock.Lock()
		if queue.status==wait{
			queue.ch<-true
			queue.status = run
		}
		queue.lock.Unlock()
	}else{
		queue.lock.RUnlock()
	}
}
func (queue* queue) poll()message.IMessage  {
	for{
		queue.lock.Lock()
		ele:=queue.messageList.Front()
		if ele!=nil {
			queue.messageList.Remove(ele)
			queue.lock.Unlock()
			return ele.Value.(message.IMessage)
		}else{
			queue.status = wait
			queue.lock.Unlock()
			<-queue.ch
		}
	}
	return nil
}

