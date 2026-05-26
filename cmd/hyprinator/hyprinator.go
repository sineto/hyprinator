package hyprinator

import "gitlab.com/sn1o/devbox/hyprinator/pkg/hyprgo"

type CLI struct {
	ipc *hyprgo.IPCClient
	event *hyprgo.EventClient
}

func New(ipc *hyprgo.IPCClient, event *hyprgo.EventClient) *CLI {
	return &CLI{ipc: ipc, event: event}
}
