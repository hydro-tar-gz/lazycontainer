```
██╗      █████╗ ███████╗██╗   ██╗ ██████╗ ██████╗ ███╗   ██╗████████╗ █████╗ ██╗███╗   ██╗███████╗██████╗ 
██║     ██╔══██╗╚══███╔╝╚██╗ ██╔╝██╔════╝██╔═══██╗████╗  ██║╚══██╔══╝██╔══██╗██║████╗  ██║██╔════╝██╔══██╗
██║     ███████║  ███╔╝  ╚████╔╝ ██║     ██║   ██║██╔██╗ ██║   ██║   ███████║██║██╔██╗ ██║█████╗  ██████╔╝
██║     ██╔══██║ ███╔╝    ╚██╔╝  ██║     ██║   ██║██║╚██╗██║   ██║   ██╔══██║██║██║╚██╗██║██╔══╝  ██╔══██╗
███████╗██║  ██║███████╗   ██║   ╚██████╗╚██████╔╝██║ ╚████║   ██║   ██║  ██║██║██║ ╚████║███████╗██║  ██║
╚══════╝╚═╝  ╚═╝╚══════╝   ╚═╝    ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝   ╚═╝   ╚═╝  ╚═╝╚═╝╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝
```                                                                                                          
# a lazy lxc & lxd tui

## prerequisites

- Go 1.22+
- lxd installed and running
- `lxc` available in `PATH`

## build & run

```bash
go install ./cmd/lc
lc
```

Alternative local build:

```bash
go build -o lc ./cmd/lc
./lc
```

## cli commands

```bash
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
- `enter`: open shell into selected instance (`bash` fallback `sh`)
- `e`: open exec command prompt and show output in right pane
- `s`: start/stop toggle
- `d`: delete selected instance (confirmation required)
- `tab`: switch between `UI` and `CLI` modes
- `t`: cycle right-pane tabs (`Info`, `Logs`, `Snapshots`) in `UI` mode
- `q`: quit

