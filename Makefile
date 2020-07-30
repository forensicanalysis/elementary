.PHONY: install-ui build build-ui build-server build-ui

install-ui:
	cd ui && yarn install
	go get -u github.com/asticode/go-astilectron-bundler/...

build-ui:
	cd ui && yarn build

build:
	go get github.com/cugu/go-resources/cmd/resources@v0.3.1
	resources -package assets -output elementary/assets/config.generated.go -trim "scripts/" scripts/scripts/* scripts/req*
	go mod tidy
	cd cmd/elementary && go build .

build-gui: install-ui build-ui
	sed 's_=/_=_g' ui/dist/index.html > tmp
	mv tmp ui/dist/index.html
	rm -rf cmd/ui/resources/app
	mv ui/dist cmd/elementary-gui/resources/app
	cp -r cmd/elementary-gui/resources/start/* cmd/elementary-gui/resources/app
	cd cmd/elementary-gui && astilectron-bundler

build-server: install-ui build-ui
	cp ui/dist cmd/elementary-server/dist
	go get -u github.com/markbates/pkger/cmd/pkger
	cd cmd/elementary-server && pkger -o assets
	cd cmd/elementary-server && go build .
