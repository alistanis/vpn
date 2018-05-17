package vpn

import (
	"errors"

	"os/user"

	"github.com/spf13/cobra"
)

// ExecuteRoot executes the root command
func ExecuteRoot() error {
	return setupCommands().Execute()
}

func rootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "vpn",
		Short: "start a vpn client connection",
	}
}

func addStartCommand(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "start [flags] [network]",
		Short: "starts a client connection on the specified network",
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("must give a network to start")
			}

			u, err := user.Current()
			if err != nil {
				return err
			}

			return start(args[0], u)

		},
	}
	root.AddCommand(cmd)
}

func addStopCommand(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "stop [flags] [network]",
		Short: "stops a client connection on the specified network",
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("must give a network to stop")
			}
			return nil
		},
	}
	root.AddCommand(cmd)
}

func setupCommands() *cobra.Command {
	root := rootCmd()
	addStartCommand(root)
	addStopCommand(root)
	return root
}
