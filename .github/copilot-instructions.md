# Copilot Instructions for dae

## Architecture & Runtime
- Entry to the daemon is in [cmd/run.go](cmd/run.go); it loads configs, wires up logging, waits for network readiness, and drives hot reload/suspend via the signal orchestration in [cmd/reload.go](cmd/reload.go) and [cmd/suspend.go](cmd/suspend.go).
- The control plane built in [control/control_plane.go](control/control_plane.go) enforces kernel feature gates, loads eBPF objects defined under [control/control.go](control/control.go) and [control/kern](control/kern), binds LAN/WAN interfaces, and keeps kernel/user routing tables in sync.
- Core services live under `component/*`: DNS orchestration, outbound dialer groups, routing, and interface monitoring (see [component/interface_manager.go](component/interface_manager.go)). Prefer extending these components instead of adding ad-hoc goroutines.

## Configuration & Routing
- Config parsing is data-driven: [config/config.go](config/config.go) declares sections with `mapstructure` tags, while [config/parser.go](config/parser.go) converts ANTLR `config_parser.Section` nodes from [pkg/config_parser](pkg/config_parser). When adding fields, update tags (`default`, `required`, `repeatable`) and ensure patches still run.
- Routing and DNS rules compile to shared kernel+userspace matchers through [component/routing](component/routing) helpers invoked inside [control/control_plane.go](control/control_plane.go). Keep new match criteria compatible with the existing optimizer/builder chain or reloads will fail verifier checks.
- Subscription and static node data flow from `subscription` and `node` sections into dialer groups before health checks; maintain `tagToNodeList` semantics so filters/policies defined in config remain predictable.

## Kernel & Networking Constraints
- `dae` is Linux-only and depends on modern kernel features (BPF loop, timers, sk_assign). The user guide in [docs/en/README.md](docs/en/README.md) documents required versions/configs; update it alongside any change that touches `consts.*FeatureVersion` or sysctl behavior.
- Respect `conf.Global.AutoConfigKernelParameter`: [control/control_plane.go](control/control_plane.go) toggles IPv4 forwarding, IPv6 `accept_ra`, and TCP fast-redirects only when that flag is set.
- Interface binding/globbing is handled by [component/interface_manager.go](component/interface_manager.go); register callbacks there instead of polling `/sys/class/net`, and reuse the manager when reacting to NIC hotplug events.

## Build & Test Workflow
- eBPF assets are regenerated via `make ebpf` (see [Makefile](Makefile)); it runs `go generate ./control/control.go` and `./trace/trace.go` with LLVM/Clang settings derived from `CLANG`, `TARGET`, and `MAX_MATCH_SET_LEN`. Always rerun before building when touching C sources or eBPF maps.
- `make dae` sets `GOOS=linux`, `CGO_ENABLED=0`, and reuses the build tags stored in `.build_tags`. Use `NOSTRIP=y` to keep symbols, or override `MAX_MATCH_SET_LEN` to stay in sync with kernel map assumptions defined in `common/consts`.
- Low-level BPF tests run through `make ebpf-test`, which regenerates fixtures and executes `go test ./control/kern/tests/...`. Standard Go formatting/testing uses `make fmt` plus `go test ./...`; some packages (e.g., control) expect root privileges or network namespaces.
- [go.mod](go.mod) pins Go 1.22 with toolchain `go1.23.2` and replaces `github.com/qimaoww/outbound`/`github.com/daeuniverse/quic-go` with sibling directories. Keep this repo checked out alongside those modules or adjust the replace directives when vendoring.

## Ops & Diagnostics
- Runtime metadata lives in `/var/run`: `dae.pid`, `dae.progress`, and `dae.abort` (see [cmd/run.go](cmd/run.go)). Any new command should reuse these paths so systemd units and scripts keep compatibility.
- Logging is centralized through `pkg/logger` invoked from [cmd/run.go](cmd/run.go); honor `Global.LogLevel`, `disable-timestamp`, and optional lumberjack rotation before creating ad-hoc loggers.
- DNS caches and eBPF objects survive reloads via `ControlPlane.EjectBpf()/InjectBpf()` and `CloneDnsCache()`; avoid storing mutable global state outside the control-plane struct or hot reloads (`dae reload`) will leak resources.
- Optional diagnostics include `conf.Global.PprofPort` (pprof listener) and the `trace` subcommand (depends on kernel ≥5.15). Document new debug surfaces inside [docs/en/README.md](docs/en/README.md) so operators can discover them quickly.
