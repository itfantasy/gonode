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

func (k *KcpNetWorker) Listen(url string) error {

	go nets.AutoPing(k)

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
			k.onError(conn, err)
			continue
		}
		go k.h_kcpSocket(conn)
	}
	return nil
}

func (k *KcpNetWorker) h_kcpSocket(conn net.Conn) {
	buf := make([]byte, 4096, 4096) // the rev buf
	defer func() {
		buf = nil
	}()
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			k.onError(conn, err)
			break
		}
		if n > 0 {
			id, exists := nets.GetInfoIdByConn(conn)
			var temp []byte = make([]byte, 0, n)
			datas := bytes.NewBuffer(temp)
			datas.Write(buf[0:n])
			if exists {
				k.onMsg(conn, id, datas.Bytes())
			} else {
				if err := k.dealHandShake(conn, string(datas.Bytes())); err != nil {
					k.onError(conn, err)
				}
			}
			temp = nil // dispose the temp buffer
		} else {
			k.onError(conn, errors.New("receive no datas!!"))
		}
	}
}

func (k *KcpNetWorker) Connect(id string, url string, origin string) error {
	theUrl := strings.Trim(url, "kcp://") // trim the ws header
	infos := strings.Split(theUrl, "/")   // parse the sub path

	conn, err := kcp.Dial(infos[0])
	if err == nil {
		k.doHandShake(conn, origin, id)
		go k.h_kcpSocket(conn)
	}
	return err
}

func (k *KcpNetWorker) Send(conn net.Conn, msg []byte) error {
	defer func() {
		msg = nil // dispose the send buffer
	}()
	_, err := conn.Write(msg)
	return err
}

func (k *KcpNetWorker) onConn(conn net.Conn, id string) {
	// record the set from id to conn
	err := nets.AddConnInfo(id, nets.KCP, conn, k)
	if err != nil {
		k.onError(conn, err)
	} else {
		k.eventListener.OnConn(id)
	}
}

func (k *KcpNetWorker) onMsg(conn net.Conn, id string, msg []byte) {
	nets.ResetConnState(id)
	if msg[0] == 35 { // '#'
		strmsg := string(msg)
		if strmsg == "#pong" {
			fmt.Println("receive pong from.." + id) // nothing to do but ResetConnState for AutoPing
			return
		} else if strmsg == "#ping" {
			fmt.Println("re sending pong to..." + id)
			go k.Send(conn, []byte("#pong")) // return the pong pck
			return
		}
	}
	k.eventListener.OnMsg(id, msg)
}

func (k *KcpNetWorker) onClose(conn net.Conn) {
	id, exists := nets.GetInfoIdByConn(conn)
	if exists {
		k.eventListener.OnClose(id)
		nets.RemoveConnInfo(id) // remove the closed conn from local record
		conn.Close()
	}
}

func (k *KcpNetWorker) onError(conn net.Conn, err error) {
	if conn != nil {
		id, exists := nets.GetInfoIdByConn(conn)
		if exists {
			k.eventListener.OnError(id, err)
			k.onClose(conn) // close the conn with errors
		} else {
			k.eventListener.OnError("", err)
			k.onClose(conn) // close the conn with errors
		}
	} else {
		k.eventListener.OnError("", err)
	}
}

func (k *KcpNetWorker) BindEventListener(eventListener nets.INetEventListener) error {
	if k.eventListener == nil {
		k.eventListener = eventListener
		return nil
	}
	return errors.New("k net worker has binded an event listener!!")
}

func (k *KcpNetWorker) Close(id string, conn net.Conn) error {
	k.eventListener.OnClose(id)
	nets.RemoveConnInfo(id) // remove the closed conn from local record
	return conn.Close()
}

func (k *KcpNetWorker) doHandShake(conn net.Conn, origin string, id string) error {
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
	k.onConn(conn, id)
	return nil
}

func (k *KcpNetWorker) dealHandShake(conn net.Conn, context string) error {
	var datas map[string]string
	if err := jsoniter.Unmarshal([]byte(context), &datas); err != nil {
		return err
	}
	origin, exists := datas["Origin"]
	if !exists {
		return errors.New("handshake datas missing!")
	}
	id, b := k.eventListener.OnCheckNode(origin) // let the gonode to check if the url is legal
	if b {
		k.onConn(conn, id)
		fmt.Println("handshake succeed !!")
		return nil
	} else {
		return errors.New("handshake illegal!!")
	}
}
