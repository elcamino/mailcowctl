package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/tob/mailcowctl/internal/client"
	"github.com/tob/mailcowctl/internal/output"
)

func newSyncjobCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{Use: "syncjob", Short: "Manage sync jobs"}
	cmd.AddCommand(syncjobListCmd(opts), syncjobGetCmd(opts), syncjobCreateCmd(opts), syncjobEditCmd(opts), syncjobDeleteCmd(opts))
	return cmd
}

func syncjobListCmd(opts *rootOptions) *cobra.Command {
	var mailbox string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List sync jobs",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := opts.client()
			if err != nil {
				return err
			}
			jobs, err := c.ListSyncJobs(opts.context(), mailbox)
			if err != nil {
				return err
			}
			return output.Render(opts.out, opts.output, jobs)
		},
	}
	cmd.Flags().StringVar(&mailbox, "mailbox", "", "filter by target mailbox")
	return cmd
}

func syncjobGetCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a sync job",
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
			job, err := c.GetSyncJob(opts.context(), id)
			if err != nil {
				return err
			}
			return output.Render(opts.out, opts.output, job)
		},
	}
}

func syncjobCreateCmd(opts *rootOptions) *cobra.Command {
	var mailbox, host1, user1, enc1, passwordEnv, exclude string
	var port1, interval, maxAge int
	var passwordStdin, delete1, delete2, delete2dupes, subscribeAll, automap bool
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a sync job",
		RunE: func(cmd *cobra.Command, args []string) error {
			if mailbox == "" || host1 == "" || user1 == "" {
				return fmt.Errorf("--mailbox, --host1, and --user1 are required")
			}
			if enc1 != "SSL" && enc1 != "TLS" && enc1 != "PLAIN" {
				return fmt.Errorf("--enc1 must be SSL, TLS, or PLAIN")
			}
			password, err := passwordFromFlags(cmd, opts.in, passwordEnv, passwordStdin)
			if err != nil {
				return err
			}
			req := client.SyncJobCreate{
				Username: mailbox, Host1: host1, Port1: port1, User1: user1, Password1: password,
				Enc1: enc1, MinsInterval: interval, MaxAge: maxAge,
				Active:            client.BoolString(true),
				Delete1:           client.BoolString(delete1),
				Delete2:           client.BoolString(delete2),
				Delete2Duplicates: client.BoolString(delete2dupes),
				SubscribeAll:      client.BoolString(subscribeAll),
				Automap:           client.BoolString(automap),
				Exclude:           exclude,
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
			return c.CreateSyncJob(opts.context(), req)
		},
	}
	cmd.Flags().StringVar(&mailbox, "mailbox", "", "local target mailbox")
	cmd.Flags().StringVar(&host1, "host1", "", "remote host")
	cmd.Flags().IntVar(&port1, "port1", 143, "remote port")
	cmd.Flags().StringVar(&user1, "user1", "", "remote username")
	cmd.Flags().StringVar(&passwordEnv, "password1-env", "", "env var holding remote password")
	cmd.Flags().BoolVar(&passwordStdin, "password1-stdin", false, "read remote password from stdin")
	cmd.Flags().StringVar(&enc1, "enc1", "TLS", "remote encryption: SSL, TLS, or PLAIN")
	cmd.Flags().IntVar(&interval, "interval", 20, "minutes between runs")
	cmd.Flags().IntVar(&maxAge, "maxage", 0, "max age in days (0 = all)")
	cmd.Flags().BoolVar(&delete1, "delete1", false, "delete from source after sync")
	cmd.Flags().BoolVar(&delete2, "delete2", false, "delete on target if missing on source")
	cmd.Flags().BoolVar(&delete2dupes, "delete2duplicates", false, "delete duplicates on target")
	cmd.Flags().BoolVar(&subscribeAll, "subscribeall", false, "subscribe all folders")
	cmd.Flags().BoolVar(&automap, "automap", false, "auto-map folders")
	cmd.Flags().StringVar(&exclude, "exclude", "", "exclude regex")
	cmd.Flags().Bool("active", false, "set active")
	cmd.Flags().Bool("inactive", false, "set inactive")
	return cmd
}

func syncjobEditCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit <id>",
		Short: "Edit a sync job",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("id must be numeric")
			}
			attrs := map[string]any{}
			stringFlagAttr(cmd, "host1", "host1", attrs)
			stringFlagAttr(cmd, "user1", "user1", attrs)
			if cmd.Flags().Changed("enc1") {
				enc1val, _ := cmd.Flags().GetString("enc1")
				if enc1val != "SSL" && enc1val != "TLS" && enc1val != "PLAIN" {
					return fmt.Errorf("--enc1 must be SSL, TLS, or PLAIN")
				}
			}
			stringFlagAttr(cmd, "enc1", "enc1", attrs)
			stringFlagAttr(cmd, "exclude", "exclude", attrs)
			intFlagAttr(cmd, "port1", "port1", attrs)
			intFlagAttr(cmd, "interval", "mins_interval", attrs)
			intFlagAttr(cmd, "maxage", "maxage", attrs)
			active, err := changedBool(cmd, "active", "inactive")
			if err != nil {
				return err
			}
			if active != nil {
				attrs["active"] = client.BoolString(*active)
			}
			if len(attrs) == 0 {
				return fmt.Errorf("no changes specified")
			}
			c, err := opts.client()
			if err != nil {
				return err
			}
			return c.EditSyncJob(opts.context(), id, attrs)
		},
	}
	cmd.Flags().String("host1", "", "remote host")
	cmd.Flags().String("user1", "", "remote username")
	cmd.Flags().String("enc1", "", "remote encryption: SSL, TLS, or PLAIN")
	cmd.Flags().String("exclude", "", "exclude regex")
	cmd.Flags().Int("port1", 0, "remote port")
	cmd.Flags().Int("interval", 0, "minutes between runs")
	cmd.Flags().Int("maxage", 0, "max age in days")
	cmd.Flags().Bool("active", false, "set active")
	cmd.Flags().Bool("inactive", false, "set inactive")
	return cmd
}

func syncjobDeleteCmd(opts *rootOptions) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a sync job",
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
			return c.DeleteSyncJob(opts.context(), id)
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "confirm deletion")
	return cmd
}
