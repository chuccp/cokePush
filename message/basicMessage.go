package message

type BasicMessage struct {
	*Message
}
func (bm *BasicMessage)GetMessageType() int8{
	return BasicMessageType
}

func createBasicMessage(fromUser string,toUser string,messageText string)* BasicMessage {
	bm:=&BasicMessage{}
	bm.SetValue(FromUser,fromUser)
	bm.SetValue(ToUser,toUser)
	bm.SetValue(MessageText,messageText)
	return bm
}