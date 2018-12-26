package tcp

import (
	"bytes"
	"errors"
	"net"
	"strings"

	"github.com/itfantasy/gonode/nets"
	"github.com/json-iterator/go"
)

type TcpNetWorker struct {
	eventListener nets.INetEventListener
}

func (this *TcpNetWorker) Listen(url string) error {
	nets.InitKvvk()
	url = strings.Trim(url, "tcp://") // trim the ws header
	infos := strings.Split(url, "/")  // parse the sub path

	tcpAddr, err := net.ResolveTCPAddr("tcp", infos[0])
	if err != nil {
		return err
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			this.onError(conn, err)
			continue
		}
		go this.h_tcpSocket(conn)
	}
	return nil
}

func (this *TcpNetWorker) h_tcpSocket(conn net.Conn) {
	var buf [4096]byte // the rev buf
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			this.onError(conn, err)
			break
		}
		if n > 0 {
			id, exists := nets.GetInfoIdByConn(conn)
			var temp []byte = make([]byte, 0, n)
			datas := bytes.NewBuffer(temp)
			datas.Write(buf[0:n])
			if exists {
				this.onMsg(conn, id, datas.Bytes())
			} else {
				if err := this.dealHandShake(conn, string(datas.Bytes())); err != nil {
					this.onError(conn, err)
				}
			}
		} else {
			this.onError(conn, errors.New("receive no datas!!"))
		}
	}
}

func (this *TcpNetWorker) Connect(url string, origin string) error {
	nets.InitKvvk()

	theUrl := strings.Trim(url, "tcp://") // trim the ws header
	infos := strings.Split(theUrl, "/")   // parse the sub path

	tcpAddr, err := net.ResolveTCPAddr("tcp", infos[0])
	if err != nil {
		return err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err == nil {
		this.doHandShake(conn, origin, url)
		go this.h_tcpSocket(conn)
	}
	return err
}

func (this *TcpNetWorker) Send(conn net.Conn, msg []byte) error {
	defer func() {
		msg = nil // dispose the send buffer
	}()
	_, err := conn.Write(msg)
	return err
}

func (this *TcpNetWorker) onConn(conn net.Conn, id string) {
	// record the set from id to conn
	err := nets.AddConnInfo(id, nets.TCP, conn)
	if err != nil {
		this.onError(conn, err)
	} else {
		this.eventListener.OnConn(id)
	}
}

func (this *TcpNetWorker) onMsg(conn net.Conn, id string, msg []byte) {
	this.eventListener.OnMsg(id, msg)
}

func (this *TcpNetWorker) onClose(conn net.Conn) {
	id, exists := nets.GetInfoIdByConn(conn)
	if exists {
		this.eventListener.OnClose(id)
		nets.RemoveConnInfo(id) // remove the closed conn from local record
		conn.Close()
	}
}

func (this *TcpNetWorker) onError(conn net.Conn, err error) {
	if conn != nil {
		id, exists := nets.GetInfoIdByConn(conn)
		if exists {
			this.eventListener.OnError(id, err)
			this.onClose(conn) // close the conn with errors
		} else {
			this.eventListener.OnError("", err)
			this.onClose(conn) // close the conn with errors
		}
	} else {
		this.eventListener.OnError("", err)
	}
}

func (this *TcpNetWorker) BindEventListener(eventListener nets.INetEventListener) error {
	if this.eventListener == nil {
		this.eventListener = eventListener
		return nil
	}
	return errors.New("this net worker has binded an event listener!!")
}

func (this *TcpNetWorker) Close(id string, conn net.Conn) error {
	this.eventListener.OnClose(id)
	nets.RemoveConnInfo(id) // remove the closed conn from local record
	return conn.Close()
}

func (this *TcpNetWorker) doHandShake(conn net.Conn, origin string, url string) error {
	info := make(map[string]string)
	info["Origin"] = origin
	datas, err := jsoniter.Marshal(info)
	if err != nil {
		return err
	}
	_, err2 := conn.Write(datas)
	if err2 != nil {
		return err2
	}
	id, legal := this.eventListener.CheckUrlLegal(url) // let the gonode to check if the url is legal
	if legal {
		this.onConn(conn, id)
		return nil
	} else {
		return errors.New("handshake illegal!!")
	}
}

func (this *TcpNetWorker) dealHandShake(conn net.Conn, info string) error {
	var datas map[string]string
	if err := jsoniter.Unmarshal([]byte(info), &datas); err != nil {
		return err
	}
	origin, exists := datas["Origin"]
	if !exists {
		return errors.New("handshake datas missing!")
	}
	url := origin
	id, legal := this.eventListener.CheckUrlLegal(url) // let the gonode to check if the url is legal
	if legal {
		this.onConn(conn, id)
		return nil
	} else {
		return errors.New("handshake illegal!!")
	}
}
