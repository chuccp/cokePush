package util

import (
	"sync"
	"time"
)

func Millisecond() uint32 {
	ms := time.Now().UnixNano() / 1e6
	return uint32(ms)
}

var TimestampFormat = "2006-01-02 15:04:05"

func FormatTime(tm *time.Time) string {
	if tm == nil {
		return ""
	}
	return tm.Format(TimestampFormat)
}

type Timer struct {
	t     *time.Timer
	isEnd chan bool
}

func newTimer() *Timer {
	return &Timer{t: time.NewTimer(time.Second * 10), isEnd: make(chan bool)}
}
func (timer *Timer) Wait() bool {
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
func (timer *Timer) End() {
	fa := timer.t.Stop()
	if fa{
		timer.isEnd <- true
	}
}
func (timer *Timer) reset(duration time.Duration) {
	timer.t.Reset(duration)
}

var poolChan = &sync.Pool{
	New: func() interface{} {
		return newTimer()
	},
}

func GetTimer(duration time.Duration) *Timer {
	timer := poolChan.Get().(*Timer)
	timer.reset(duration)
	return timer
}
func FreeTimer(timer *Timer) {
	poolChan.Put(timer)
}
