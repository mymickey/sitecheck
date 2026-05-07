
# SiteCheck

SiteCheck is a macOS menubar-only network connectivity monitor built with Wails v3, Vue 3, Pinia, and official HeroUI styles.

# 起因

灵感来源于 [IPCheck.ing](https://ipcheck.ing) 的网络探测方式, 致敬阿禅

# roadmap
- [ ] 折线图统计过去24小时的延迟
- [x] 增加 My IP 地址显示
- [ ] 系统主题自适应
- [ ] 还有很多细节需要打磨

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

## License

This software is for **personal and non-commercial use only**. If you intend to use this software for commercial purposes (including but not limited to use within a business, company, or for-profit entity), you must obtain a commercial license.

For commercial licensing inquiries or further information, please contact the author.


