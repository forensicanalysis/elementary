.PHONY: build-cli build-ui build-server build-ui

build-ui:
	cd ui && yarn install
	cd ui && yarn build

build-cli:
	go mod tidy
	cd cmd/elementary && go build .

build-gui: build-ui
	sed 's_="/_="_g' ui/dist/index.html > tmp
	mv tmp ui/dist/index.html
	rm -rf cmd/elementary-gui/resources/app
	mv ui/dist cmd/elementary-gui/resources/app
	cp -r cmd/elementary-gui/resources/start/* cmd/elementary-gui/resources/app
	rm -f cmd/elementary-gui/bind.go
	go get -u github.com/asticode/go-astilectron-bundler/...
	rm -f cmd/elementary-gui/bindata.dummy.go
	cd cmd/elementary-gui && astilectron-bundler

build-server: build-ui
	rm -rf cmd/elementary-server/dist
	cp -r ui/dist cmd/elementary-server/dist
	go mod tidy
	cd cmd/elementary-server && go build .
