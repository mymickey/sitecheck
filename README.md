# SiteCheck

SiteCheck is a macOS menubar-only network connectivity monitor built with Wails v3, Vue 3, Pinia, and official HeroUI styles.

## Development

Use the Makefile so the project always switches to Go 1.26.1 before starting Wails:

```sh
make dev
```

Other useful commands:

```sh
make test
make build
make dmg
```

`make dmg` creates `bin/dmg/SiteCheck.dmg`.

## Behavior

- The app runs as an accessory app with no Dock icon.
- The menubar item renders the app logo plus the fastest and slowest latency from the latest benchmark.
- The native menu contains `Connectivity`, `Setting`, and `Quit`.
- `Connectivity > Benchmark` probes the five configured targets and updates both the native menu labels and menubar summary.
- The settings window manages the five targets and the background benchmark interval, defaulting to 10 minutes.

The probe method mirrors IPCheck.ing Connectivity: GET the target URL with `Cache-Control: no-store`, treat any HTTP response as reachable, and mark DNS/connect/timeout failures as unavailable.


### roadmap
- [ ] 折线图统计过去24小时的延迟
- [ ] 增加 My IP 地址显示
