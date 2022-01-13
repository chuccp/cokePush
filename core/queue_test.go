package core

import (
	log "github.com/chuccp/coke-log"
	"os"
	"os/signal"
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
	root   *element
	ch     chan bool
	isWait bool
}

func newQueue() *queue {
	return &queue{ch: make(chan bool)}
}
func (queue *queue) Offer(value interface{}) {
	ele := newElement(value)
	if queue.root == nil {
		queue.root = ele
	} else {
		queue.root.next = ele
		ele.prev = queue.root
	}
	if queue.isWait {
		queue.ch <- true
	}
}
func (queue *queue) Poll() interface{} {
	for {
		if queue.root != nil {
			value := queue.root.value
			if queue.root.next != nil {
				queue.root = queue.root.next
				queue.root.prev = nil
			} else {
				queue.root = nil
			}
			return value
		} else {
			queue.wait()
		}
	}
}
func (queue *queue) wait() {
	queue.isWait = true
	<-queue.ch
}

func TestCompare(t *testing.T) {

	que := newQueue()

	go func() {
		for {
			v := que.Poll()
			log.Info(v)
		}
	}()

	time.Sleep(time.Second)

	que.Offer(1)
	time.Sleep(time.Second)

	que.Offer(2)
	que.Offer(3)
	time.Sleep(time.Second)

	que.Offer(4)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGBUS)
	<-sig
}
