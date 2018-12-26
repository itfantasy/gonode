package nets

import "net"

type INetWorker interface {
	Listen(url string) error
	Connect(url string, origin string) error
	Send(conn net.Conn, msg []byte) error
	BindEventListener(eventListener INetEventListener) error
	Close(id string, conn net.Conn) error
}

const (
	None string = "none"
	WS   string = "ws"
	TCP  string = "tcp"
	KCP  string = "kcp"
	HTTP string = "http"
)
