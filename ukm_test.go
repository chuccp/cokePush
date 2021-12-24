package main

import (
	log "github.com/chuccp/coke-log"
	"github.com/chuccp/cokePush/ukm"
	"testing"
)

func Test_oo(t *testing.T)  {
	listen, err := ukm.ListenKm(7676)
	if err == nil {
		for {
			conn, err := listen.AcceptConn()
			if err == nil {
				go func() {
					data := make([]byte, 1024)
					for {
						num, err := conn.Read(data)
						if err == nil {
							log.Info(conn.RemoteAddr().String(), ";", string(data[:num]))
						}
					}
				}()
			} else {
				log.Info("{}",err)
			}
		}
	}
}
