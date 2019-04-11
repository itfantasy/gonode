package nets

import (
	"errors"
	"net"
)

type INetWorker interface {
	Listen(url string) error
	Connect(id string, url string, origin string) error
	Send(conn net.Conn, msg []byte) error
	Close(id string, conn net.Conn) error
	BindEventListener(eventListener INetEventListener) error
}

const (
	None string = "none"
	WS   string = "ws"
	TCP  string = "tcp"
	KCP  string = "kcp"
	HTTP string = "http"
)

func LocalIPv4s() ([]string, error) {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.IP.String())
		}
	}

	return ips, nil
}

func IsLocalIPv4(ip string) error {
	ips, err := LocalIPv4s()
	if err != nil {
		return err
	}
	for _, item := range ips {
		if item == ip {
			return nil
		}
	}
	return errors.New("illegal local IPv4!!!")
}
