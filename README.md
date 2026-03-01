```
██╗      █████╗ ███████╗██╗   ██╗ ██████╗ ██████╗ ███╗   ██╗████████╗ █████╗ ██╗███╗   ██╗███████╗██████╗ 
██║     ██╔══██╗╚══███╔╝╚██╗ ██╔╝██╔════╝██╔═══██╗████╗  ██║╚══██╔══╝██╔══██╗██║████╗  ██║██╔════╝██╔══██╗
██║     ███████║  ███╔╝  ╚████╔╝ ██║     ██║   ██║██╔██╗ ██║   ██║   ███████║██║██╔██╗ ██║█████╗  ██████╔╝
██║     ██╔══██║ ███╔╝    ╚██╔╝  ██║     ██║   ██║██║╚██╗██║   ██║   ██╔══██║██║██║╚██╗██║██╔══╝  ██╔══██╗
███████╗██║  ██║███████╗   ██║   ╚██████╗╚██████╔╝██║ ╚████║   ██║   ██║  ██║██║██║ ╚████║███████╗██║  ██║
╚══════╝╚═╝  ╚═╝╚══════╝   ╚═╝    ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝   ╚═╝   ╚═╝  ╚═╝╚═╝╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝
```                                                                                                          
# a lazy lxc & lxd tui

## a fair warning before you use this

- this may not work on arch please dont use this on arch unless you know how to configure it

## prerequisites

- Go 1.22+
- lxd installed and running
- `lxc` available in `PATH`

## build & run

if `lc` is not found after `go install` (which it probably isn't), add go's bin dir to your shell `PATH` first.

bash:

```
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc
source ~/.bashrc
```

zsh:

```
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc
source ~/.zshrc
```

fish:

```
fish_add_path (go env GOPATH)/bin
```

then install and run:

```
go install ./cmd/lc
lc
```

alternative local build:

```
go build -o lc ./cmd/lc
./lc
```

## cli commands

```
lc ls [--json]
lc up <image> --name <name>
lc start <name>
lc stop <name>
lc rm <name> [-f]
lc shell <name>
lc exec <name> -- <cmd...>
lc ip <name>
lc snap <name> <snapshot>
lc restore <name> <snapshot>
```

## tui keybinds

- `j` / `k` or arrows: move selection
- `/`: focus search filter
- `r`: refresh list
- `u`: create container with image picker/autofill 
- `enter`: open shell into selected instance 
- `e`: open exec command prompt and show output in right pane
- `s`: start/stop toggle
- `d`: delete selected instance 
- `tab`: switch between `UI` and `CLI` modes
- `t`: cycle right-pane tabs
- `q`: quit
