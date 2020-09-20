.PHONY: build-cli build-ui build-server build-ui

build-ui:
	cd ui && yarn install
	cd ui && yarn build

build-cli:
	go get -u github.com/asticode/go-bindata/...
	rm -f cmd/elementary/bindata.dummy.go
	go-bindata -prefix "scripts/" -o cmd/elementary/bindata.generated.go scripts/...
	go mod tidy
	cd cmd/elementary && go build .

build-gui: build-ui
	sed 's_=/_=_g' ui/dist/index.html > tmp
	mv tmp ui/dist/index.html
	rm -rf cmd/ui/resources/app
	mv ui/dist cmd/elementary-gui/resources/app
	cp -r cmd/elementary-gui/resources/start/* cmd/elementary-gui/resources/app
	rm -f cmd/elementary-gui/bind.go
	go get -u github.com/asticode/go-astilectron-bundler/...
	rm cmd/elementary-gui/bindata.dummy.go
	cd cmd/elementary-gui && astilectron-bundler

build-server: build-ui
	go get -u github.com/asticode/go-bindata/...
	rm cmd/elementary-server/bindata.dummy.go
	go-bindata -prefix "ui/dist/" -o cmd/elementary-server/bindata.generated.go ui/dist/...
	go mod tidy
	cd cmd/elementary-server && go build .
