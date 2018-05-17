package vpn

import (
	"errors"

	"os/user"

	"strings"

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
	var mfa bool
	aliases := []string{"up"}
	cmd := &cobra.Command{
		Use:     "start [flags] [network]",
		Short:   "starts a client connection on the specified network. Aliases \"" + strings.Join(aliases, " ") + "\"",
		Aliases: aliases,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("must give a network to start")
			}

			u, err := user.Current()
			if err != nil {
				return err
			}

			return start(args[0], u, mfa)

		},
	}

	cmd.Flags().BoolVarP(&mfa, "mfa", "m", false, "will search for a secret in ~/.totp/secret for use with mfa vpns")

	root.AddCommand(cmd)
}

func addStopCommand(root *cobra.Command) {
	aliases := []string{"down", "kill"}
	cmd := &cobra.Command{
		Use:     "stop [flags] [network]",
		Short:   "stops a client connection on the specified network. Aliases \"" + strings.Join(aliases, " ") + "\"",
		Aliases: aliases,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("must give a network to stop")
			}
			return stop(args[0])
		},
	}
	root.AddCommand(cmd)
}

func addStatusCommand(root *cobra.Command) {
	aliases := []string{"ls", "stat"}
	cmd := &cobra.Command{
		Use:     "status [network]",
		Short:   "reports the status of the vpn on the given network. Aliases \"" + strings.Join(aliases, ", ") + "\"",
		Aliases: aliases,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) > 0 {
				return status(args[0])
			}
			return status("")
		},
	}
	root.AddCommand(cmd)
}

func addTotpCommand(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "totp",
		Short: "generates a totp password for use with a mfa vpn",
		RunE: func(c *cobra.Command, args []string) error {
			u, err := user.Current()
			if err != nil {
				return err
			}
			return otp(u)
		},
	}
	root.AddCommand(cmd)
}

func addRestartCommand(root *cobra.Command) {
	var mfa bool
	cmd := &cobra.Command{
		Use:   "restart [network]",
		Short: "restarts a vpn on the given network",
		RunE: func(c *cobra.Command, args []string) error {
			// we don't care about this error
			stop(args[0])
			u, err := user.Current()
			if err != nil {
				return err
			}
			return start(args[0], u, mfa)
		},
		Args: cobra.MinimumNArgs(1),
	}

	cmd.Flags().BoolVarP(&mfa, "mfa", "m", false, "will search for a secret in ~/.totp/secret for use with mfa vpns")
	root.AddCommand(cmd)
}

func setupCommands() *cobra.Command {
	root := rootCmd()
	addStartCommand(root)
	addStopCommand(root)
	addStatusCommand(root)
	addRestartCommand(root)
	addTotpCommand(root)
	return root
}
