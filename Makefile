.PHONY: install-ui build build-ui build-server build-ui pack-server

install-ui:
	cd ui && yarn install
	go get -u github.com/asticode/go-astilectron-bundler/...

build-ui:
	cd ui && yarn build

pack-cli:
	go get github.com/cugu/go-resources/cmd/resources@v0.3.1
	resources -package assets -output cmd/elementary/assets/config.generated.go -trim "scripts/" scripts/scripts/* scripts/req*
	go mod tidy

build: pack-cli
	cd cmd/elementary && go build .

pack-gui: install-ui build-ui
	sed 's_=/_=_g' ui/dist/index.html > tmp
	mv tmp ui/dist/index.html
	rm -rf cmd/ui/resources/app
	mv ui/dist cmd/elementary-gui/resources/app
	cp -r cmd/elementary-gui/resources/start/* cmd/elementary-gui/resources/app

build-gui: pack-gui
	cd cmd/elementary-gui && astilectron-bundler

pack-server: install-ui build-ui
	cp -r ui/dist cmd/elementary-server/dist
	go get -u github.com/markbates/pkger/cmd/pkger
	mkdir cmd/elementary-server/assets
	pkger -o cmd/elementary-server/assets

build-server: pack-server
	cd cmd/elementary-server && go build .
