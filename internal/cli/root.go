package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"lazycontainer/internal/core"
	"lazycontainer/internal/tui"
)

func NewRootCmd(backend core.Backend) *cobra.Command {
	root := &cobra.Command{
		Use:           "lc",
		Aliases:       []string{"lazy"},
		Short:         "lazycontainer: lazydocker-like UI + CLI for LXD",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return tui.Run(backend)
		},
	}
	root.PersistentFlags().Bool("no-color", false, "Disable colorized output")

	root.AddCommand(newUICmd(backend))
	root.AddCommand(newLsCmd(backend))
	root.AddCommand(newUpCmd(backend))
	root.AddCommand(newStartCmd(backend))
	root.AddCommand(newStopCmd(backend))
	root.AddCommand(newRmCmd(backend))
	root.AddCommand(newShellCmd(backend))
	root.AddCommand(newExecCmd(backend))
	root.AddCommand(newIPCmd(backend))
	root.AddCommand(newSnapCmd(backend))
	root.AddCommand(newRestoreCmd(backend))

	return root
}

func newUICmd(backend core.Backend) *cobra.Command {
	return &cobra.Command{
		Use:   "ui",
		Short: "Start terminal UI",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tui.Run(backend)
		},
	}
}

func newLsCmd(backend core.Backend) *cobra.Command {
	var jsonOut bool
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List LXD instances",
		RunE: func(cmd *cobra.Command, args []string) error {
			instances, err := backend.ListInstances(context.Background())
			if err != nil {
				return err
			}
			p := newPalette(cmd, cmd.OutOrStdout())
			if jsonOut {
				b, err := json.MarshalIndent(instances, "", "  ")
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(b))
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\t%s\t%s\n", p.Header("NAME"), p.Header("STATE"), p.Header("IP"), p.Header("IMAGE"))
			if len(instances) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), p.Muted("(no instances)"))
				return nil
			}
			for _, inst := range instances {
				fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\t%s\t%s\n", p.Name(inst.Name), p.State(inst.State), inst.IP, inst.Image)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&jsonOut, "json", false, "Output as JSON")
	return cmd
}

func newUpCmd(backend core.Backend) *cobra.Command {
	var name string
	cmd := &cobra.Command{
		Use:   "up <image> --name <name>",
		Short: "Create and start an instance from image",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(name) == "" {
				return errors.New("--name is required")
			}
			return backend.Launch(context.Background(), args[0], name)
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "Instance name")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newStartCmd(backend core.Backend) *cobra.Command {
	return &cobra.Command{
		Use:   "start <name>",
		Short: "Start an instance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return backend.Start(context.Background(), args[0])
		},
	}
}

func newStopCmd(backend core.Backend) *cobra.Command {
	return &cobra.Command{
		Use:   "stop <name>",
		Short: "Stop an instance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return backend.Stop(context.Background(), args[0])
		},
	}
}

func newRmCmd(backend core.Backend) *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "rm <name>",
		Short: "Delete an instance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return backend.Delete(context.Background(), args[0], force)
		},
	}
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force delete")
	return cmd
}

func newShellCmd(backend core.Backend) *cobra.Command {
	return &cobra.Command{
		Use:   "shell <name>",
		Short: "Open an interactive shell in instance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return backend.Shell(context.Background(), args[0])
		},
	}
}

func newExecCmd(backend core.Backend) *cobra.Command {
	return &cobra.Command{
		Use:   "exec <name> -- <cmd...>",
		Short: "Run command in instance",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			out, err := backend.Exec(context.Background(), name, args[1:])
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), out)
			return nil
		},
	}
}

func newIPCmd(backend core.Backend) *cobra.Command {
	return &cobra.Command{
		Use:   "ip <name>",
		Short: "Get instance IP",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ip, err := backend.IP(context.Background(), args[0])
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), ip)
			return nil
		},
	}
}

func newSnapCmd(backend core.Backend) *cobra.Command {
	return &cobra.Command{
		Use:   "snap <name> <snapshot>",
		Short: "Create snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return backend.Snapshot(context.Background(), args[0], args[1])
		},
	}
}

func newRestoreCmd(backend core.Backend) *cobra.Command {
	return &cobra.Command{
		Use:   "restore <name> <snapshot>",
		Short: "Restore snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return backend.Restore(context.Background(), args[0], args[1])
		},
	}
}
