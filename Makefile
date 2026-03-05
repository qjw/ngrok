.PHONY: server client release-all

BUILDTAGS=debug

server:
	go build -tags '$(BUILDTAGS)' -o build/ngrokd-$(BUILDTAGS) ./cmd/ngrokd

client:
	go build -tags '$(BUILDTAGS)' -o build/ngrok-$(BUILDTAGS) ./cmd/ngrok

release-client: BUILDTAGS=release
release-client: client

release-server: BUILDTAGS=release
release-server: server

release-all: release-client release-server