package util

import (
	log "github.com/chuccp/coke-log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func Test_timer2(t *testing.T) {
	cd := make(chan bool, 2)

	go func() {

		log.Info(<-cd)

	}()
	time.Sleep(time.Second * 10)
	for i := 0; i <= 1000; i++ {
		tm := time.NewTimer(time.Second * 10)
		go func() {
			<-tm.C
			cd <- true
			log.Info("111")
		}()
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGBUS)
	<-sig
}
func Test_timer(t *testing.T) {

	var queue = NewQueue()

	go func() {
		for {
			v, num := queue.Poll()
			log.Info(v, "====", num)
		}
	}()
	for i := 0; i <= 5_000_000; i++ {
		tm := time.NewTimer(time.Second * 10)
		go func() {
			<-tm.C
			queue.Offer(123)
		}()
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGBUS)
	<-sig

}
