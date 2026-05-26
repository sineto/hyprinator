package hyprgo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"unicode"
)

// Dispatch wrap 'hyprctl dispatch' command
func (ipc *IPCClient) Dispatch(params ...string) ([]byte, error) {
	raw, err := ipc.composeCommand("dispatch", params, false)
	if err != nil {
		return nil, err
	}

	resp, err := ipc.request(raw)
	if err != nil {
		return nil, err
	}

	body := strings.TrimRight(string(resp), "\n\r")
	if body != SUCCESS {
		return nil, fmt.Errorf("hyprland dispatch failed: %s", body)
	}

	return []byte(body), nil
}

// Clients wrap 'hyprclt clients' command
func (ipc *IPCClient) Clients(params ...string) ([]Client, error) {
	raw, err := ipc.composeCommand("clients", params, true)
	if err != nil {
		return nil, err
	}

	resp, err := ipc.request(raw)
	if err != nil {
		return nil, err
	}

	var clients []Client
	if err := json.Unmarshal(resp, &clients); err != nil {
		return nil, fmt.Errorf("parse clients: %w", err)
	}

	return clients, nil
}

func (ev *EventClient) WaitForOpenWindow(ctx context.Context)(OpenWindow, error) {
	for {
		msg, err := ev.Receive(ctx)
		if err != nil {
			return OpenWindow{}, err
		}

		for _, data := range msg {
			if data.Type == EventOpenWindow {
				rawEvent := strings.Split(data.Data, ",")
				if len(rawEvent) > 3 {
					return OpenWindow{
						Address: rawEvent[0],
						Class: rawEvent[1],
						Title: rawEvent[2],
					}, nil
				}
			}
		}
	}
}

func (ipc *IPCClient) composeCommand(command string, params []string, jflag bool) ([]byte, error) {
	if command == "" {
		panic("empty command")
	}

	buf := bytes.NewBuffer(nil)

	if jflag {
		buf.Write([]byte{'j', '/'})
	}

	buf.WriteString(command)
	buf.WriteString(" ")
	for _, param := range params {
		buf.WriteString(param)
		buf.WriteString(" ")
	}

	// _ = json.NewEncoder(os.Stdin).Encode(string(bytes.TrimRightFunc(buf.Bytes(), unicode.IsSpace)))

	return bytes.TrimRightFunc(buf.Bytes(), unicode.IsSpace), nil
}
