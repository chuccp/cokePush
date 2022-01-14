package util

import (
	"sync"
	"sync/atomic"
	"time"
)

type element struct {
	next  *element
	value interface{}
}

func newElement(value interface{}) *element {
	return &element{value: value}
}

type Queue struct {
	input  *element
	output *element
	ch     chan bool
	isWait bool
	num    int32
	rLock  *sync.RWMutex
	timer  *timer
}

type timer struct {
	t     *time.Timer
	isEnd chan bool
}

func newTimer() *timer {
	return &timer{t: time.NewTimer(time.Second * 10), isEnd: make(chan bool)}
}
func (timer *timer) wait() bool {
	select {
	case <-timer.t.C:
		{
			return false
		}
	case fa := <-timer.isEnd:
		{
			return fa
		}
	}
	return false
}
func (timer *timer) end() {
	fa := timer.t.Stop()
	if fa {
		timer.isEnd <- true
	}
}
func (timer *timer) reset(duration time.Duration) {
	timer.t.Reset(duration)
}

var poolTimer = &sync.Pool{
	New: func() interface{} {
		return newTimer()
	},
}

func getTimer(duration time.Duration) *timer {
	ti := poolTimer.Get().(*timer)
	ti.reset(duration)
	return ti
}
func freeTimer(timer *timer) {
	poolTimer.Put(timer)
}

func NewQueue() *Queue {
	return &Queue{ch: make(chan bool), isWait: false, num: 0, rLock: new(sync.RWMutex)}
}
func (queue *Queue) Offer(value interface{}) (num int32) {
	ele := newElement(value)
	queue.rLock.Lock()
	if queue.num == 0 {
		queue.input = ele
		queue.output = ele
	} else {
		queue.input.next = ele
		queue.input = ele
	}
	num = atomic.AddInt32(&queue.num, 1)
	if queue.isWait == true {
		queue.isWait = false
		queue.rLock.Unlock()
		queue.ch <- true
	} else {
		queue.rLock.Unlock()
	}
	return
}
func (queue *Queue) Num()int32{
	return queue.num
}
func (queue *Queue) Poll() (value interface{}, num int32) {
	for {
		queue.rLock.Lock()
		if queue.num > 0 {
			if queue.num == 1 {
				var ele = queue.output
				queue.num--
				num = queue.num
				queue.rLock.Unlock()
				return ele.value, num
			} else {
				queue.rLock.Unlock()
				var ele = queue.output
				value = ele.value
				queue.output = ele.next
				ele.next = nil
				num = atomic.AddInt32(&queue.num, -1)
				return
			}
		} else {
			queue.isWait = true
			queue.rLock.Unlock()
			<-queue.ch
		}
	}
}

func (queue *Queue) Take(duration time.Duration) (value interface{}, num int32) {
	for {
		queue.rLock.Lock()
		if queue.num > 0 {
			if queue.num == 1 {
				var ele = queue.output
				queue.num--
				num = queue.num
				queue.rLock.Unlock()
				return ele.value, num
			} else {
				queue.rLock.Unlock()
				var ele = queue.output
				value = ele.value
				queue.output = ele.next
				ele.next = nil
				num = atomic.AddInt32(&queue.num, -1)
				return
			}
		} else {
			queue.isWait = true
			queue.rLock.Unlock()
			queue.timer = getTimer(duration)
			go func() {
				fa := queue.timer.wait()
				if !fa {
					queue.rLock.Lock()
					if queue.isWait == true {
						queue.isWait = false
						queue.rLock.Unlock()
						queue.ch <- false
					}
				}
			}()
			flag := <-queue.ch
			queue.timer.end()
			freeTimer(queue.timer)
			if !flag {
				return nil, 0
			}
		}
	}
}


