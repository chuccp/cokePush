package clog

import (
	log "github.com/chuccp/coke-log"
)

func Default()*log.Config  {
	config:=log.GetConfig()
	config.SetLevel(log.TraceLevel)
	return config
}
