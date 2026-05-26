package hyprsinator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"gitlab.com/sn1o/devbox/hyprinator/pkg/hyprgo"
)

type SpaceItem struct {
	App       string   `yaml:"app"`
	Exec      []string `yaml:"exec"`
	Cwd       string   `yaml:"cwd"`
	Workspace int      `yaml:"workspace"`
	Monitor   string   `yaml:"monitor"`
	Focus     bool     `yaml:"focus"`
	Delay     int      `yaml:"delay"`
}

func (cli *CLI) SetupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup <space>",
		Short: "Will setup <space> environment config",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, _ := cmd.Flags().GetString("config")
			return cli.Setup(args[0], configPath)
		},
	}

	cmd.Flags().StringP("config", "c", "$HOME/.hyprinator.yaml", "Path to config file")
	return cmd
}

func (cli *CLI) Setup(spaceName, configPath string) error {
	path, err := expandPath(configPath)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	var config map[string][]SpaceItem
	if err = yaml.Unmarshal(data, &config); err != nil {
		return err
	}

	items, ok := config[spaceName]
	if !ok {
		return fmt.Errorf("space %q not found in hs.yaml", spaceName)
	}

	var focusedItem SpaceItem
	var focusedItemAddress string
	for _, item := range items {
		_, err = cli.ipc.Dispatch("workspace", strconv.Itoa(item.Workspace))
		if err != nil {
			return err
		}

		_, err = cli.ipc.Dispatch("moveworkspacetomonitor", strings.Join([]string{strconv.Itoa(item.Workspace), item.Monitor}, " "))
		if err != nil {
			return err
		}

		args := []string{strconv.Quote(item.App)}

		if item.Cwd != "" {
			args = append(args, "--working-directory", strconv.Quote(item.Cwd))
		}

		if len(item.Exec) > 0 {
			execStr := strings.Join(item.Exec, " ; ")
			args = append(args, "-e", "sh","-c", strconv.Quote(execStr))
		}

		command := strings.Join(args, " ")
		_, err = cli.ipc.Dispatch("exec", command)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		win, err := cli.event.WaitForOpenWindow(ctx)
		cancel()

		if err != nil {
			return err
		}

		if item.Focus {
			focusedItem = item
			focusedItemAddress = win.Address
		}

		time.Sleep(time.Duration(item.Delay) * time.Millisecond)
	}

	_, err = cli.ipc.Dispatch("workspace", strconv.Itoa(focusedItem.Workspace))
	if err != nil {
		return err
	}

	_, err = cli.ipc.Dispatch("focuswindow", "address:0x"+focusedItemAddress)
	if err != nil {
		return err
	}

	return nil
}

func (cli *CLI) ExecAndWaitOpenWindow(command string, timeout time.Duration) (hyprgo.OpenWindow, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	flushCtx, flushCancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer flushCancel()

	_, err := cli.event.Receive(flushCtx)
	if err != nil {
		return hyprgo.OpenWindow{}, err
	}

	_, err = cli.ipc.Dispatch("exec", command)
	if err != nil {
		return hyprgo.OpenWindow{}, err
	}

	win, err := cli.event.WaitForOpenWindow(ctx)
	if err != nil {
		return hyprgo.OpenWindow{}, err
	}

	return win, nil
}

func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		path = strings.Replace(path, "~", home, 1)
	}

	path = os.ExpandEnv(path)

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return filepath.Clean(absPath), nil
}
