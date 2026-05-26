package hyprgo

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
)

func NewIPCClient(socket string) *IPCClient {
	return &IPCClient{
		conn: &net.UnixAddr{
			Net:  "unix",
			Name: socket,
		},
	}
}

func (ipc *IPCClient) request(request []byte) ([]byte, error) {
	if len(request) == 0 {
		return nil, fmt.Errorf("empty request")
	}

	conn, err := net.Dial("unix", ipc.conn.Name)
	if err != nil {
		return nil, fmt.Errorf("connect to Hyprland socket: %w", err)
	}
	defer conn.Close()

	writer := bufio.NewWriter(conn)

	_, err = writer.Write(request)
	if err != nil {
		return nil, fmt.Errorf("write to Hyprland socket: %w", err)
	}

	err = writer.Flush()
	if err != nil {
		return nil, fmt.Errorf("error while flushing to socket: %w", err)
	}

	// Get the response back
	rbuf := bytes.NewBuffer(nil)
	sbuf := make([]byte, BUFSIZE)
	reader := bufio.NewReader(conn)

	for {
		n, er := reader.Read(sbuf)
		if er != nil {
			if er == io.EOF {
				break
			}

			return nil, fmt.Errorf("error while reading from socket: %w", er)
		}

		rbuf.Write(sbuf[:n])

		if n < BUFSIZE {
			break
		}
	}

	return rbuf.Bytes(), err
}
