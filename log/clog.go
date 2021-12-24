package clog

import log "github.com/chuccp/coke-log"

func Start()  {
	config:=log.GetConfig()
	config.SetLevel(log.InfoLevel)
}
