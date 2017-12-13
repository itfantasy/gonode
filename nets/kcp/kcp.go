package kcp

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/itfantasy/gonode/nets"
	"github.com/xtaci/kcp-go"
)

type KcpNetWorker struct {
	eventListener nets.INetEventListener
}

func (this *KcpNetWorker) Listen(url string) error {
	this.initKvvk()
	url = strings.Trim(url, "kcp://") // trim the ws header
	infos := strings.Split(url, "/")  // parse the sub path
	listener, err := kcp.Listen(infos[0])
	if err != nil {
		return err
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			if conn != nil {
				this.onError(conn, err)
			} else {
				fmt.Println(err) // an err without conn
			}
			continue
		}
		go this.h_kcpSocket(conn)
	}
	return nil
}

func (this *KcpNetWorker) h_kcpSocket(conn net.Conn) {
	url := conn.RemoteAddr().String()
	id, legal := this.eventListener.CheckUrlLegal(url) // let the gonode to check if the url is legal
	if legal {
		this.onConn(conn, id)
		var buf [4096]byte // the rev buf
		for {
			n, err := conn.Read(buf[0:])
			if err != nil {
				this.onError(conn, err)
				break
			}
			if n > 0 {
				datas := bytes.NewBuffer(nil)
				datas.Write(buf[0:n])
				this.onMsg(conn, datas.Bytes())
			}
		}
	} else {
		conn.Close() // dispose the illegel conn
	}
}

func (this *KcpNetWorker) Connect(url string, origin string) error {
	this.initKvvk()

	url = strings.Trim(url, "kcp://") // trim the ws header
	infos := strings.Split(url, "/")  // parse the sub path

	conn, err := kcp.Dial(infos[0])
	if err == nil {
		go this.h_kcpSocket(conn)
	}
	return err
}

func (this *KcpNetWorker) Send(id string, msg []byte) error {
	defer func() {
		msg = nil // dispose the send buffer
	}()
	conn, exist := this.getInfoConnById(id)
	if exist {
		_, err := conn.Write(msg)
		return err
	} else {
		return errors.New("there is not the conn for this id in local record!")
	}
}

func (this *KcpNetWorker) SendAsync(id string, msg []byte) {
	go this.Send(id, msg)
}

func (this *KcpNetWorker) onConn(conn net.Conn, id string) {
	// record the set from id to conn
	err := this.addConnInfo(id, conn)
	if err != nil {
		this.onError(conn, err)
	} else {
		this.eventListener.OnConn(id)
	}
}

func (this *KcpNetWorker) onMsg(conn net.Conn, msg []byte) {
	id, exists := this.getInfoIdByConn(conn)
	if exists {
		this.eventListener.OnMsg(id, msg)
	}
}

func (this *KcpNetWorker) onClose(conn net.Conn) {
	id, exists := this.getInfoIdByConn(conn)
	if exists {
		this.eventListener.OnClose(id)
		this.removeConnInfo(id) // remove the closed conn from local record
		conn.Close()
	}
}

func (this *KcpNetWorker) onError(conn net.Conn, err error) {
	id, exists := this.getInfoIdByConn(conn)
	if exists {
		this.eventListener.OnError(id, err)
		this.onClose(conn) // close the conn with errors
	}
}

func (this *KcpNetWorker) BindEventListener(eventListener nets.INetEventListener) error {
	if this.eventListener == nil {
		this.eventListener = eventListener
		return nil
	}
	return errors.New("this net worker has binded an event listener!!")
}

func (this *KcpNetWorker) Close(id string) error {
	conn, exists := this.getInfoConnById(id)
	if exists {
		this.eventListener.OnClose(id)
		this.removeConnInfo(id) // remove the closed conn from local record
		return conn.Close()
	}
	return errors.New("there is not the id in local record!")
}
