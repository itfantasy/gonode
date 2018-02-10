package kcp

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/itfantasy/gonode/nets"
	"github.com/json-iterator/go"
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
			this.onError(conn, err)
			continue
		}
		go this.h_kcpSocket(conn)
	}
	return nil
}

func (this *KcpNetWorker) h_kcpSocket(conn net.Conn) {
	var buf [4096]byte // the rev buf
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			this.onError(conn, err)
			break
		}
		if n > 0 {
			id, exists := this.getInfoIdByConn(conn)
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

func (this *KcpNetWorker) Connect(url string, origin string) error {
	this.initKvvk()

	theUrl := strings.Trim(url, "kcp://") // trim the ws header
	infos := strings.Split(theUrl, "/")   // parse the sub path

	conn, err := kcp.Dial(infos[0])
	if err == nil {
		this.doHandShake(conn, origin, url)
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

func (this *KcpNetWorker) onMsg(conn net.Conn, id string, msg []byte) {
	this.eventListener.OnMsg(id, msg)
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
	if conn != nil {
		id, exists := this.getInfoIdByConn(conn)
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

func (this *KcpNetWorker) doHandShake(conn net.Conn, origin string, url string) error {
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

func (this *KcpNetWorker) dealHandShake(conn net.Conn, info string) error {
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
		fmt.Println("handshake succeed !!")
		return nil
	} else {
		return errors.New("handshake illegal!!")
	}
}

func (this *KcpNetWorker) keepAlive() {

}
