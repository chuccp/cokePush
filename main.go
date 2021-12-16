package main

import (
	"github.com/chuccp/cokePush/api"
	"github.com/chuccp/cokePush/config"
	"github.com/chuccp/cokePush/core"
	"github.com/chuccp/cokePush/ukm"
	"log"
)

func DefaultRegister() *core.Register {
	config := config.DefaultConfig()
	var defaultRegister = core.NewRegister()
	defaultRegister.AddServer(api.NewServer(config))
	return defaultRegister
}
func main2() {
	reg := DefaultRegister()
	cokePush := reg.Create()
	cokePush.StartSync()
}
func main() {
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
							log.Println(conn.RemoteAddr().String(), ";", string(data[:num]))
						}
					}
				}()
			} else {
				log.Println(err)
			}

		}
	}
}
