package assets

import (
	_ "embed"
)

var (
	//go:embed client/tls/ngrokroot.crt
	CaCrt []byte

	//go:embed server/tls/server.crt
	ServerCrt []byte
	//go:embed server/tls/server.key
	ServerKey []byte

	//go:embed client/tls/client.crt
	ClientCrt []byte
	//go:embed client/tls/client.key
	ClientKey []byte
)
