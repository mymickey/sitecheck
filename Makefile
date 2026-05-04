SHELL := /bin/zsh
WAILS_VITE_PORT ?= 9245

.PHONY: sync-icons dev build test package dmg

sync-icons:
	cp frontend/public/sitecheck-app.png build/appicon.png
	rm -f build/darwin/Assets.car
	source "$$HOME/.gvm/scripts/gvm" && gvm use go1.26.1 && wails3 generate icons -input build/appicon.png -macfilename build/darwin/icons.icns -windowsfilename build/windows/icon.ico

dev: sync-icons
	source "$$HOME/.gvm/scripts/gvm" && gvm use go1.26.1 && wails3 dev -config ./build/config.yml -port $(WAILS_VITE_PORT)

# `make build`: compile the latest binary only (`bin/sitecheck`)
build: sync-icons
	source "$$HOME/.gvm/scripts/gvm" && gvm use go1.26.1 && wails3 task build

test:
	source "$$HOME/.gvm/scripts/gvm" && gvm use go1.26.1 && go test ./... && npm --prefix frontend run build

# `make package`: build the macOS `.app` bundle only (`bin/sitecheck.app`)
package: build
	rm -rf bin/sitecheck.app bin/sitecheck.dev.app
	source "$$HOME/.gvm/scripts/gvm" && gvm use go1.26.1 && wails3 task package

# `make dmg`: the normal release/install flow; rebuild everything and create `bin/dmg/SiteCheck.dmg`
dmg: package
	mkdir -p bin/dmg
	rm -rf bin/dmg-root
	mkdir -p bin/dmg-root
	cp -R bin/sitecheck.app bin/dmg-root/SiteCheck.app
	ln -s /Applications bin/dmg-root/Applications
	hdiutil create -volname SiteCheck -srcfolder bin/dmg-root -ov -format UDZO bin/dmg/SiteCheck.dmg
