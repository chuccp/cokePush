package message

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