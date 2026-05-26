package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"gitlab.com/sn1o/devbox/hyprinator/cmd/hyprinator"
	"gitlab.com/sn1o/devbox/hyprinator/pkg/hyprgo"
	"gitlab.com/sn1o/devbox/hyprinator/pkg/hyprgo/helpers"
)


func RootCmd(cli *hyprinator.CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hyprinator",
		Short: "hyprinator - tmuxinator-like setup workspaces environment in Hyprland",
	}

	cmd.AddCommand(cli.SetupCommand())
	return cmd
}

func Execute() {
	wsocket, err := helpers.GetSocketPath(hyprgo.RequestSocket)
	if err != nil {
		os.Exit(1)
	}

	rsocket, err := helpers.GetSocketPath(hyprgo.EventSocket)
	if err != nil {
		os.Exit(1)
	}

	ipc := hyprgo.NewIPCClient(wsocket)
	event := hyprgo.NewEventClient(rsocket)

	cli := hyprinator.New(ipc, event)

	if err := RootCmd(cli).Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error while running command '%s'", err)
		os.Exit(1)
	}
}
