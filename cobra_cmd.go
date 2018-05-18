package vpn

import (
	"errors"

	"os/user"

	"strings"

	"fmt"

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

			return start(args[0], u)

		},
	}
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
			return otp(u, args[0])
		},
		Args: cobra.MinimumNArgs(1),
	}
	root.AddCommand(cmd)
}

func addRestartCommand(root *cobra.Command) {
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
			return start(args[0], u)
		},
		Args: cobra.MinimumNArgs(1),
	}

	root.AddCommand(cmd)
}

type config struct {
	URL       string `json:"url"`
	MFASecret string `json:"mfa_secret"`
	Port      int    `json:"port"`
}

// URLString returns a URLString suitable for use in openvpn commands
func (c *config) URLString() string {
	return fmt.Sprintf("%s %d", c.URL, c.Port)
}

func addConfigureCommand(root *cobra.Command) {
	c := &config{}
	cmd := &cobra.Command{
		Use:   "configure [flags] [network]",
		Short: "creates the necessary configuration files for vpn + totp - ca.crt, user.key, and user.crt must be provided separately",
		Long: `
configre will create configuration files for this vpn.

vpn expects user certs/keys to live in $HOMEDIR/.vpn/$VPNUSERNAME-$NETWORK, so /Users/ccooper/.vpn/ccooper-iad or $HOMEDIR/.vpn/ccooper-default

this command will create a config.json file in each network directory.

configure will create these directories if they do not exist, but it is the user's responsibility to place the correct ca.crt, username.key, and username.crt in each directory

to be clear, username.key and username.crt should be your actual vpn username, so in my case they would be ccooper.key and ccooper.crt

ca.crt should be called ca.crt exactly
`,
		RunE: func(command *cobra.Command, args []string) error {
			u, err := user.Current()
			if err != nil {
				return err
			}
			return configure(args[0], c, u)
		},
		Args: cobra.MinimumNArgs(1),
	}

	cmd.Flags().StringVarP(&c.URL, "url", "u", "", "the url for this remote")
	cmd.Flags().IntVarP(&c.Port, "port", "p", -1, "the port for this remote")
	cmd.Flags().StringVarP(&c.MFASecret, "mfa-secret", "m", "", "the mfa secret key to generate one time passwords")

	root.AddCommand(cmd)
}

func setupCommands() *cobra.Command {
	root := rootCmd()
	addStartCommand(root)
	addStopCommand(root)
	addStatusCommand(root)
	addRestartCommand(root)
	addTotpCommand(root)
	addConfigureCommand(root)
	return root
}
