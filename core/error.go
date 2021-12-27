package core

import "errors"

var NoFoundUser = errors.New("NOT FOUND USER")

var UnKnownMessageType = errors.New("UNKNOWN MESSAGE TYPE")

var UnKnownClassId = errors.New("UNKNOWN CLASS ID")

var MessageIdIsBlank = errors.New("message Id is blank")

var  ProtocolError = errors.New("km protocolError")


var  ReadTimeout = errors.New("read timeout")
