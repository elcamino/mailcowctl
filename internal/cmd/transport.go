package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/tob/mailcowctl/internal/client"
	"github.com/tob/mailcowctl/internal/output"
)

func newTransportCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{Use: "transport", Short: "Manage outbound transport maps"}
	cmd.AddCommand(transportListCmd(opts), transportGetCmd(opts), transportCreateCmd(opts), transportDeleteCmd(opts))
	return cmd
}

func transportListCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List transport maps",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := opts.client()
			if err != nil {
				return err
			}
			list, err := c.ListTransports(opts.context())
			if err != nil {
				return err
			}
			return output.Render(opts.out, opts.output, list)
		},
	}
}

func transportGetCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a transport map",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("id must be numeric")
			}
			c, err := opts.client()
			if err != nil {
				return err
			}
			t, err := c.GetTransport(opts.context(), id)
			if err != nil {
				return err
			}
			return output.Render(opts.out, opts.output, t)
		},
	}
}

func transportCreateCmd(opts *rootOptions) *cobra.Command {
	var destination, nexthop, username, passwordEnv string
	var passwordStdin bool
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a transport map",
		RunE: func(cmd *cobra.Command, args []string) error {
			if destination == "" || nexthop == "" {
				return fmt.Errorf("--destination and --nexthop are required")
			}
			password := ""
			if passwordEnv != "" || passwordStdin {
				var err error
				password, err = passwordFromFlags(cmd, opts.in, passwordEnv, passwordStdin)
				if err != nil {
					return err
				}
			}
			req := client.TransportCreate{
				Destination: destination, Nexthop: nexthop, Username: username, Password: password,
				Active: client.BoolString(true),
			}
			active, err := changedBool(cmd, "active", "inactive")
			if err != nil {
				return err
			}
			if active != nil {
				req.Active = client.BoolString(*active)
			}
			c, err := opts.client()
			if err != nil {
				return err
			}
			return c.CreateTransport(opts.context(), req)
		},
	}
	cmd.Flags().StringVar(&destination, "destination", "", "destination domain or address")
	cmd.Flags().StringVar(&nexthop, "nexthop", "", "next hop, e.g. smtp:[relay.example.org]:587")
	cmd.Flags().StringVar(&username, "username", "", "auth username (optional)")
	cmd.Flags().StringVar(&passwordEnv, "password-env", "", "env var holding the auth password")
	cmd.Flags().BoolVar(&passwordStdin, "password-stdin", false, "read the auth password from stdin")
	cmd.Flags().Bool("active", false, "set active")
	cmd.Flags().Bool("inactive", false, "set inactive")
	return cmd
}

func transportDeleteCmd(opts *rootOptions) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a transport map",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("id must be numeric")
			}
			if err := confirm(opts.in, opts.errOut, yes, args[0]); err != nil {
				return err
			}
			c, err := opts.client()
			if err != nil {
				return err
			}
			return c.DeleteTransport(opts.context(), id)
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "confirm deletion")
	return cmd
}
