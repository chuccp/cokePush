package message

type Write interface {
	WriteMessage(iMessage IMessage) error
}
