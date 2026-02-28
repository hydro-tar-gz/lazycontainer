package main

import (
	"fmt"
	"os"

	"lazycontainer/internal/backend/lxd"
	"lazycontainer/internal/cli"
)

func main() {
	backend := lxd.NewClient(nil)
	root := cli.NewRootCmd(backend)
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
