package util

import (
	log "github.com/chuccp/coke-log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func Test_queue(t *testing.T)  {

	var q = NewQueue()

	for i := 0; i < 5; i++ {
		go func() {
			j := i
			for  {

				v,_:=q.Poll()
				time.Sleep(time.Second*1)
				log.Info(j,"=====",v)
			}

		}()
		time.Sleep(time.Second)
	}

	time.Sleep(time.Second*10)

	q.Offer(1)
	q.Offer(2)
	q.Offer(3)
	q.Offer(4)
	q.Offer(5)

	q.Offer(1)
	q.Offer(2)
	q.Offer(3)
	q.Offer(4)
	q.Offer(5)

	q.Offer(1)
	q.Offer(2)
	q.Offer(3)
	q.Offer(4)
	q.Offer(5)

	q.Offer(1)
	q.Offer(2)
	q.Offer(3)
	q.Offer(4)
	q.Offer(5)

	q.Offer(1)
	q.Offer(2)
	q.Offer(3)
	q.Offer(4)
	q.Offer(5)

	q.Offer(1)
	q.Offer(2)
	q.Offer(3)
	q.Offer(4)
	q.Offer(5)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGBUS)
	<-sig

}
