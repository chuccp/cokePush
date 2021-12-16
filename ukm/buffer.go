package ukm

import (
	"bytes"
	"io"
	"sync"
)

type buffer struct {
	buff *bytes.Buffer
	r    chan bool
	wait bool
	lock *sync.RWMutex
}

func (buffer *buffer) Write(data []byte) (int, error) {
	buffer.lock.Lock()
	defer buffer.lock.Unlock()
	num, err := buffer.buff.Write(data)
	if buffer.wait {
		buffer.r <- true
		buffer.wait = false
	}
	return num, err
}
func (buffer *buffer) Read(p []byte) (int, error) {
	buffer.lock.RLock()
	num, err := buffer.buff.Read(p)
	if err == io.EOF {
		buffer.wait = true
		buffer.lock.RUnlock()
		<-buffer.r
		buffer.lock.RLock()
		num, err := buffer.buff.Read(p)
		buffer.lock.RUnlock()
		return num, err
	}
	buffer.lock.RUnlock()
	return num, nil
}

func newBuffer(data []byte) *buffer {
	var bu = bytes.NewBuffer(data)
	return &buffer{buff: bu, r: make(chan bool), wait: false, lock: new(sync.RWMutex)}
}
