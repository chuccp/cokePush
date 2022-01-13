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
	next  *element
	value interface{}
}

func newElement(value interface{}) *element {
	return &element{value: value}
}

type queue struct {
	input  *element
	output *element
	ch     chan bool
	isWait int32
	num    int32
	rLock  *sync.RWMutex
}

func newQueue() *queue {
	return &queue{ch: make(chan bool), isWait: 0, num: 0, rLock: new(sync.RWMutex)}
}
func (queue *queue) Offer(value interface{}) (num int32) {
	ele := newElement(value)
	queue.rLock.Lock()
	if atomic.CompareAndSwapInt32(&(queue.num), 0, 1) {
		queue.input = ele
		queue.output = ele
		num = 1
		queue.rLock.Unlock()
	} else {
		queue.input.next = ele
		queue.input = ele
		num = atomic.AddInt32(&queue.num, 1)
		queue.rLock.Unlock()
	}
	if atomic.CompareAndSwapInt32(&(queue.isWait), 1, 0) {
		queue.ch <- true
	}
	return
}
func (queue *queue) Poll() (value interface{}, num int32) {
	for {
		v := atomic.LoadInt32(&queue.num)
		if v > 0 {
			queue.rLock.Lock()
			if atomic.CompareAndSwapInt32(&(queue.num), 1, 0) {
				var ele = queue.output
				value = ele.value
				queue.output = ele.next
				ele.next = nil
				queue.rLock.Unlock()
				return
			} else {
				queue.rLock.Unlock()
				var ele = queue.output
				value = ele.value
				queue.output = ele.next
				ele.next = nil
				num = atomic.AddInt32(&queue.num, -1)
			}
			return
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
	for i := 0; i < 1000; i++ {

		go func() {
			num := que.Offer(3)
			log.Info("__!!!__", num)
			num = que.Offer(3)
			log.Info("__!!!__", num)
			num = que.Offer(3)
			log.Info("__!!!__", num)
			num = que.Offer(3)
			log.Info("__!!!__", num)
			num = que.Offer(3)
			log.Info("__!!!__", num)
			num = que.Offer(3)
			log.Info("__!!!__", num)
			num = que.Offer(3)
			log.Info("__!!!__", num)
			num = que.Offer(3)
			log.Info("__!!!__", num)
			num = que.Offer(3)
			log.Info("__!!!__", num)
			num = que.Offer(3)
			log.Info("__!!!__", num)
			num = que.Offer(3)
			log.Info("__!!!__", num)
			num = que.Offer(3)
			log.Info("__!!!__", num)
			num = que.Offer(3)
			log.Info("__!!!__", num)
			num = que.Offer(3)
			log.Info("__!!!__", num)
			num = que.Offer(3)
			log.Info("__!!!__", num)

		}()

	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGBUS)
	<-sig
}
