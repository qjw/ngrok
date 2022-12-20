.PHONY: server client deps release-all assets client-assets server-assets

BUILDTAGS=debug

deps:
	go install github.com/jteeuwen/go-bindata/go-bindata@latest

server: assets
	go build -tags '$(BUILDTAGS)' -o build/ngrokd-$(BUILDTAGS) ./cmd/ngrokd

client: assets
	go build -tags '$(BUILDTAGS)' -o build/ngrok-$(BUILDTAGS) ./cmd/ngrok

assets: client-assets server-assets

client-assets:
	go-bindata -nomemcopy -pkg=assets -tags=$(BUILDTAGS) \
		-debug=$(if $(findstring debug,$(BUILDTAGS)),true,false) \
		-o=ngrok/client/assets/assets_$(BUILDTAGS).go \
		assets/client/...

server-assets:
	go-bindata -nomemcopy -pkg=assets -tags=$(BUILDTAGS) \
		-debug=$(if $(findstring debug,$(BUILDTAGS)),true,false) \
		-o=ngrok/server/assets/assets_$(BUILDTAGS).go \
		assets/server/...

release-client: BUILDTAGS=release
release-client: client

release-server: BUILDTAGS=release
release-server: server

release-all: release-client release-server