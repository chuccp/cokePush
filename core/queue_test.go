package core

import (
	log "github.com/chuccp/coke-log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"
)

type element struct {
	next, prev *element
	value      interface{}
}

func newElement(value interface{}) *element {
	return &element{value: value}
}

type queue struct {
	input  *element
	output *element
	ch     chan bool
	isWait int32
	isHas  int32
	num    int32
	rLock  *sync.RWMutex
}

func newQueue() *queue {
	return &queue{ch: make(chan bool), isWait: 0, num: 0, isHas: 0, rLock: new(sync.RWMutex)}
}
func (queue *queue) Offer(value interface{}) int32 {
	ele := newElement(value)
	if atomic.CompareAndSwapInt32(&(queue.isHas), 0, 1) {
		queue.input = ele
		queue.output = ele
	} else {
		queue.rLock.Lock()
		queue.input.next = ele
		//ele.prev = queue.input
		queue.input = ele
		queue.rLock.Unlock()
	}
	num := atomic.AddInt32(&queue.num, 1)
	if atomic.CompareAndSwapInt32(&(queue.isWait), 1, 0) {
		queue.ch <- true
	}
	return num
}
func (queue *queue) Poll() (interface{}, int32) {
	for {
		v := atomic.LoadInt32(&queue.isHas)
		if v == 1 {
			var ele = queue.output
			queue.output = ele.next
			ele.next = nil
			num := atomic.AddInt32(&queue.num, -1)
			if num == 0 {
				atomic.StoreInt32(&queue.isHas, 0)
			}
			//else{
			//	queue.output.prev = nil
			//}
			return ele.value, num
		} else {
			atomic.StoreInt32(&(queue.isWait), 1)
			queue.wait()
		}
	}
}
func (queue *queue) wait() {
	<-queue.ch
}

func TestCompare(t *testing.T) {
	que := newQueue()
	go func() {
		for {
			v, num := que.Poll()
			runtime.Gosched()
			log.Info(v, "____", num)
		}
	}()
	time.Sleep(time.Second)
	for i:=0;i<1000;i++{

		go func() {
			num:=que.Offer(3)
			log.Info("__!!!__", num)
			num=que.Offer(3)
			log.Info("__!!!__", num)
			num=que.Offer(3)
			log.Info("__!!!__", num)
			num=que.Offer(3)
			log.Info("__!!!__", num)
			num=que.Offer(3)
			log.Info("__!!!__", num)
			num=que.Offer(3)
			log.Info("__!!!__", num)
			num=que.Offer(3)
			log.Info("__!!!__", num)
			num=que.Offer(3)
			log.Info("__!!!__", num)
			num=que.Offer(3)
			log.Info("__!!!__", num)
			num=que.Offer(3)
			log.Info("__!!!__", num)
			num=que.Offer(3)
			log.Info("__!!!__", num)
			num=que.Offer(3)
			log.Info("__!!!__", num)
			num=que.Offer(3)
			log.Info("__!!!__", num)
			num=que.Offer(3)
			log.Info("__!!!__", num)
			num=que.Offer(3)
			log.Info("__!!!__", num)

		}()

	}




	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGBUS)
	<-sig
}
