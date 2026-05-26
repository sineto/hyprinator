package hyprsinator

// func (cli *CLI) DispatchCommand() *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "dispatch <dispatcher> [args]",
// 		Short: "Hyprland distpacher command",
// 		Args:  cobra.MinimumNArgs(1),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			arg := ""
// 			if len(args) > 1 {
// 				arg = strings.Join(args[1:], " ")
// 			}
//
// 			_, err := cli.ipc.Dispatch(args[0], arg)
// 			return err
// 		},
// 	}
// }
