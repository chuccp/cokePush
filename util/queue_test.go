package util

import (
	log "github.com/chuccp/coke-log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)
func Test_chan2(t *testing.T){

	flag:=make(chan bool)

	go func() {

		flag<-false
		log.Info("1")
		flag<-true
		log.Info("2")
	}()

	log.Info(<-flag)
	time.Sleep(time.Second*10)
	log.Info(<-flag)

}
func Test_queue(t *testing.T)  {

	var q = NewQueue()

	for i := 0; i < 5; i++ {
		go func() {
			j := i
			for  {

				v,_:=q.Take(time.Second)
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

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGBUS)
	<-sig

}
