package km

var  TimeoutError = &timeoutError{}

type timeoutError struct {
}

func (te timeoutError)Error() string  {
	return "io timeout"
}