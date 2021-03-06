package message

import "bytes"

type BasicMessage struct {
	*Message
}

func CreateBasicMessage(fromUser string, toUser string, messageText string) *BasicMessage {
	bm := &BasicMessage{Message: CreateMessage(OrdinaryMessageClass, BasicMessageType)}
	bm.SetString(FromUser, fromUser)
	bm.SetString(ToUser, toUser)
	bm.SetString(Text, messageText)
	return bm
}
type BackBasicMessage struct {
	*Message
}
func CreateBackBasicMessage(isSuccess bool,msgId uint32)*BackBasicMessage  {
	if isSuccess{
		bm := &BackBasicMessage{Message: CreateBackMessage(BackMessageClass, BackMessageOKType,msgId)}
		return bm
	}else{
		bm := &BackBasicMessage{Message: CreateBackMessage(BackMessageClass, BackMessageErrorType,msgId)}
		return bm
	}
}


func (basic *BasicMessage)SetExMsgId(msgId string)  {
	basic.SetString(ExMessageId, msgId)
}
func (basic *BasicMessage)GetExMsgId()string  {
	return basic.GetString(ExMessageId)
}

type MultiMessage struct {
	*Message
}

func CreateMultiMessage(fromUser string, toUser *[]string, messageText string)*MultiMessage{
	bm := &MultiMessage{Message: CreateMessage(OrdinaryMessageClass, MultiMessageType)}
	bm.SetString(FromUser, fromUser)
	var buffer  = bytes.NewBuffer([]byte{})
	for _,v:=range *toUser{
		buffer.WriteString(v)
		buffer.WriteString(";")
	}
	bm.SetValue(ToUser, buffer.Bytes())
	bm.SetString(Text, messageText)
	return bm
}

type LoginMessage struct {
	*Message
}
func CreateLoginMessage(username string)*LoginMessage  {
	bm := &LoginMessage{Message: CreateMessage(FunctionMessageClass, LoginType)}
	bm.SetString(Username, username)
	return bm
}

type LiveMessage struct {
	*Message
}

func CreateLiveMessage()*LiveMessage  {
	bm := &LiveMessage{Message: CreateMessage(LiveMessageClass, BlankLiveType)}
	return bm
}