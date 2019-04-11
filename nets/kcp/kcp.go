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
	nets.InitKvvk()
	go nets.AutoPing(this)

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
			temp = nil // dispose the temp buffer
		} else {
			this.onError(conn, errors.New("receive no datas!!"))
		}
	}
}

func (this *KcpNetWorker) Connect(id string, url string, origin string) error {
	nets.InitKvvk()

	theUrl := strings.Trim(url, "kcp://") // trim the ws header
	infos := strings.Split(theUrl, "/")   // parse the sub path

	conn, err := kcp.Dial(infos[0])
	if err == nil {
		this.doHandShake(conn, origin, url, id)
		go this.h_kcpSocket(conn)
	}
	return err
}

func (this *KcpNetWorker) Send(conn net.Conn, msg []byte) error {
	defer func() {
		msg = nil // dispose the send buffer
	}()
	_, err := conn.Write(msg)
	return err
}

func (this *KcpNetWorker) onConn(conn net.Conn, id string) {
	// record the set from id to conn
	err := nets.AddConnInfo(id, nets.KCP, conn, this)
	if err != nil {
		this.onError(conn, err)
	} else {
		this.eventListener.OnConn(id)
	}
}

func (this *KcpNetWorker) onMsg(conn net.Conn, id string, msg []byte) {
	nets.ResetConnState(id)
	if msg[0] == 35 { // '#'
		strmsg := string(msg)
		if strmsg == "#pong" {
			fmt.Println("receive pong from.." + id)
			return
		} else if strmsg == "#ping" {
			fmt.Println("re sending pong to..." + id)
			go this.Send(conn, []byte("#pong")) // return the pong pck
			return
		}
	}
	this.eventListener.OnMsg(id, msg)
}

func (this *KcpNetWorker) onClose(conn net.Conn) {
	id, exists := nets.GetInfoIdByConn(conn)
	if exists {
		this.eventListener.OnClose(id)
		nets.RemoveConnInfo(id) // remove the closed conn from local record
		conn.Close()
	}
}

func (this *KcpNetWorker) onError(conn net.Conn, err error) {
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

func (this *KcpNetWorker) BindEventListener(eventListener nets.INetEventListener) error {
	if this.eventListener == nil {
		this.eventListener = eventListener
		return nil
	}
	return errors.New("this net worker has binded an event listener!!")
}

func (this *KcpNetWorker) Close(id string, conn net.Conn) error {
	this.eventListener.OnClose(id)
	nets.RemoveConnInfo(id) // remove the closed conn from local record
	return conn.Close()
}

func (this *KcpNetWorker) doHandShake(conn net.Conn, origin string, url string, id string) error {
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

	id, b := this.eventListener.OnCheckNode(id, url) // let the gonode to check if the url is legal
	if b {
		this.onConn(conn, id)
		return nil
	} else {
		return errors.New("handshake illegal!! " + url + "#" + id)
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
	urlAndId := strings.Split(origin, "#")
	if len(urlAndId) != 2 {
		return errors.New("illegal origin data! " + origin)
	}
	id := urlAndId[1]
	url := urlAndId[0]
	id, b := this.eventListener.OnCheckNode(id, url) // let the gonode to check if the url is legal
	if b {
		this.onConn(conn, id)
		fmt.Println("handshake succeed !!")
		return nil
	} else {
		return errors.New("handshake illegal!!")
	}
}
