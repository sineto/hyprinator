package hyprgo

import "net"

const (
	BUFSIZE = 8192
	SUCCESS = "ok"

	RequestSocket = ".socket.sock"
	EventSocket = ".socket2.sock"
)

type IPCClient struct {
	conn *net.UnixAddr
}

type IPC interface {
	Dispatch(params ...string) ([]byte, error)
}

var _ IPC = (*IPCClient)(nil)

type Client struct {
	Address string `json:"address"`
	Class   string `json:"class"`
	Title   string `json:"title"`
}
