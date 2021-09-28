package main

import (
	"github.com/chuccp/cokePush/config"
	"github.com/magiconair/properties"
	"log"
)

func main() {

	cfg, err := config.LoadFile("application.properties", properties.UTF8)
	if err==nil{
		log.Println(cfg.ReadString("a"))
	}else{
		log.Println(err)
	}
}
