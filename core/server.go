package core

type Server interface {
	Start()error
	Init()
	Name()string
}


