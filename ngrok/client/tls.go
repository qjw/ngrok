package client

import (
	_ "crypto/sha512"
	"crypto/tls"
	"crypto/x509"

	"github.com/qjw/ngrok/assets"
)

func LoadTLSConfig() (*tls.Config, error) {
	// 信任CA根证书
	pool := x509.NewCertPool()
	ok := pool.AppendCertsFromPEM(assets.CaCrt)
	if !ok {
		panic("failed to parse root certificate")
	}

	// 服务端证书配置
	var (
		cert tls.Certificate
		err  error
	)

	if cert, err = tls.X509KeyPair(assets.ClientCrt, assets.ClientKey); err != nil {
		return nil, err
	}

	return &tls.Config{
		RootCAs:      pool,
		Certificates: []tls.Certificate{cert},
	}, nil
}
