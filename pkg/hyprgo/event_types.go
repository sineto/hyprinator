package hyprgo

import "net"

type EventHandler interface {
	OpenWindow(ow OpenWindow)
}

type OpenWindow struct {
	Address string
	Class   string
	Title   string
}

type EventClient struct{
	conn net.Conn
}

type EventType string

type ReceiveData struct {
	Type EventType
	Data string
}

const (
	EventOpenWindow EventType = "openwindow"
)

func (e EventClient) GetAllEvents() []EventType {
	return []EventType{
		EventOpenWindow,
	}
}
