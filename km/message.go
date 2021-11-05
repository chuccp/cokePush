package km

type Message interface {
	GetMessageType() (int, error)
}
