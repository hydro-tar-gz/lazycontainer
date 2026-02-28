package core

import "context"

type Instance struct {
	Name  string `json:"name"`
	State string `json:"state"`
	IP    string `json:"ip"`
	Image string `json:"image"`
}

type Backend interface {
	ListInstances(ctx context.Context) ([]Instance, error)
	Launch(ctx context.Context, image, name string) error
	Start(ctx context.Context, name string) error
	Stop(ctx context.Context, name string) error
	Delete(ctx context.Context, name string, force bool) error
	Exec(ctx context.Context, name string, cmd []string) (string, error)
	Shell(ctx context.Context, name string) error
	IP(ctx context.Context, name string) (string, error)
	Snapshot(ctx context.Context, name, snapshot string) error
	Restore(ctx context.Context, name, snapshot string) error
	Logs(ctx context.Context, name string) (string, error)
	Snapshots(ctx context.Context, name string) (string, error)
	ImageAliases(ctx context.Context) ([]string, error)
}
