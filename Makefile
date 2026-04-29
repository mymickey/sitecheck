SHELL := /bin/zsh
WAILS_VITE_PORT ?= 9245

.PHONY: dev build test package dmg

dev:
	source "$$HOME/.gvm/scripts/gvm" && gvm use go1.26.1 && wails3 dev -config ./build/config.yml -port $(WAILS_VITE_PORT)

build:
	source "$$HOME/.gvm/scripts/gvm" && gvm use go1.26.1 && wails3 task build

test:
	source "$$HOME/.gvm/scripts/gvm" && gvm use go1.26.1 && go test ./... && npm --prefix frontend run build

package:
	source "$$HOME/.gvm/scripts/gvm" && gvm use go1.26.1 && wails3 task package

dmg: package
	mkdir -p bin/dmg
	hdiutil create -volname SiteCheck -srcfolder bin/sitecheck.app -ov -format UDZO bin/dmg/SiteCheck.dmg
