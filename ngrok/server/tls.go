package server

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/qjw/ngrok/assets"
)

func LoadTLSConfig() (tlsConfig *tls.Config, err error) {
	// 信任CA根证书
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(assets.CaCrt)
	if !ok {
		panic("failed to parse root certificate")
	}

	// 服务端证书配置
	var cert tls.Certificate
	if cert, err = tls.X509KeyPair(assets.ServerCrt, assets.ServerKey); err != nil {
		return
	}

	tlsConfig = &tls.Config{
		ClientCAs:          clientCertPool,
		Certificates:       []tls.Certificate{cert},
		ClientAuth:         tls.RequireAndVerifyClientCert,
		InsecureSkipVerify: false,
	}

	return
}
