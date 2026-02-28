package lxd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"lazycontainer/internal/core"
)

type Client struct {
	runner Runner
}

func NewClient(r Runner) *Client {
	if r == nil {
		r = CommandRunner{}
	}
	return &Client{runner: r}
}

func (c *Client) ListInstances(ctx context.Context) ([]core.Instance, error) {
	out, err := c.runner.Run(ctx, "lxc", "list", "--format", "json")
	if err != nil {
		return nil, err
	}
	return parseListJSON(out)
}

func (c *Client) Launch(ctx context.Context, image, name string) error {
	_, err := c.runner.Run(ctx, "lxc", "launch", image, name)
	if err == nil {
		return nil
	}
	logs, logErr := c.runner.Run(ctx, "lxc", "info", name, "--show-log")
	if logErr != nil || strings.TrimSpace(string(logs)) == "" {
		return err
	}
	return fmt.Errorf("%w\n--- lxc info %s --show-log ---\n%s", err, name, strings.TrimSpace(string(logs)))
}

func (c *Client) Start(ctx context.Context, name string) error {
	_, err := c.runner.Run(ctx, "lxc", "start", name)
	return err
}

func (c *Client) Stop(ctx context.Context, name string) error {
	_, err := c.runner.Run(ctx, "lxc", "stop", name)
	return err
}

func (c *Client) Delete(ctx context.Context, name string, force bool) error {
	args := []string{"delete", name}
	if force {
		args = append(args, "--force")
	}
	_, err := c.runner.Run(ctx, "lxc", args...)
	return err
}

func (c *Client) Exec(ctx context.Context, name string, cmd []string) (string, error) {
	if len(cmd) == 0 {
		return "", fmt.Errorf("command cannot be empty")
	}
	args := []string{"exec", name, "--"}
	args = append(args, cmd...)
	out, err := c.runner.Run(ctx, "lxc", args...)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(out), "\n"), nil
}

func (c *Client) Shell(ctx context.Context, name string) error {
	if err := runInteractive(ctx, "lxc", "exec", name, "--", "bash"); err == nil {
		return nil
	}
	return runInteractive(ctx, "lxc", "exec", name, "--", "sh")
}

func runInteractive(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *Client) IP(ctx context.Context, name string) (string, error) {
	instances, err := c.ListInstances(ctx)
	if err != nil {
		return "", err
	}
	for _, inst := range instances {
		if inst.Name == name {
			if inst.IP == "" {
				return "-", nil
			}
			return inst.IP, nil
		}
	}
	return "", fmt.Errorf("instance %q not found", name)
}

func (c *Client) Snapshot(ctx context.Context, name, snapshot string) error {
	_, err := c.runner.Run(ctx, "lxc", "snapshot", name, snapshot)
	return err
}

func (c *Client) Restore(ctx context.Context, name, snapshot string) error {
	_, err := c.runner.Run(ctx, "lxc", "restore", name, snapshot)
	return err
}

func (c *Client) Logs(ctx context.Context, name string) (string, error) {
	out, err := c.runner.Run(ctx, "lxc", "info", name, "--show-log")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (c *Client) Snapshots(ctx context.Context, name string) (string, error) {
	out, err := c.runner.Run(ctx, "lxc", "info", name)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (c *Client) ImageAliases(ctx context.Context) ([]string, error) {
	out, err := c.runner.Run(ctx, "lxc", "image", "list", "images:", "--format", "json")
	if err != nil {
		return nil, err
	}
	return parseImageAliasesJSON(out)
}
