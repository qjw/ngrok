package proto

import (
	"github.com/qjw/ngrok/ngrok/conn"
)

type Protocol interface {
	GetName() string
	WrapConn(conn.Conn, interface{}) conn.Conn
}
