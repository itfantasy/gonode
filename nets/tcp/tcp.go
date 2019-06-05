package tcp

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/itfantasy/gonode/core/binbuf"
	"github.com/itfantasy/gonode/nets"
	"github.com/json-iterator/go"
)

type TcpNetWorker struct {
	eventListener nets.INetEventListener
}

func (t *TcpNetWorker) Listen(url string) error {

	go nets.AutoPing(t)

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
			t.onError(conn, err)
			continue
		}
		go t.h_tcpSocket(conn)
	}
	return nil
}

func (t *TcpNetWorker) h_tcpSocket(conn net.Conn) {
	buf := make([]byte, 8192, 8192)
	rcvbuf := NewTcpBuffer(buf)
	defer func() {
		rcvbuf.Dispose()
	}()
	for {
		n, err := conn.Read(rcvbuf.Buffer())
		if err != nil {
			t.onError(conn, err)
			break
		}
		if n > 0 {
			rcvbuf.AddDataLen(n)
			for rcvbuf.Count() > PCK_MIN_SIZE {
				parser := binbuf.BuildParser(rcvbuf.Slice(), 0)
				head, err := parser.Int()
				if err != nil || head != PCK_HEADER {
					rcvbuf.Clear()
					break
				}
				l, err := parser.Short()
				length := int(l)
				if err != nil {
					rcvbuf.Clear()
					break
				} else if length > rcvbuf.Count() {
					break
				}
				id, exists := nets.GetInfoIdByConn(conn)
				src := rcvbuf.Slice()[PCK_MIN_SIZE : PCK_MIN_SIZE+length]
				var temp []byte = make([]byte, 0, length)
				datas := bytes.NewBuffer(temp)
				datas.Write(src)
				rcvbuf.DeleteData(length + PCK_MIN_SIZE)
				if exists {
					t.onMsg(conn, id, datas.Bytes())
				} else {
					if err := t.dealHandShake(conn, string(datas.Bytes())); err != nil {
						t.onError(conn, err)
						return
					}
				}
				temp = nil // dispose the temp buffer
			}
			rcvbuf.Reset()
		} else {
			t.onError(conn, errors.New("receive no datas!!"))
		}
	}

	// or you can use the bufio.Scanner like t
	/*
		scanner := bufio.NewScanner(conn)
		scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			if len(data) > PCK_MIN_SIZE {
				if !atEOF {
					parser := binbuf.BuildParser(data, 0)
					header, err := parser.Int()
					if err == nil && header == PCK_HEADER {
						length, err2 := parser.Short()
						if err2 == nil {
							needlen := int(length) + PCK_MIN_SIZE
							if needlen <= len(data) { // parser a whole package
								return needlen, data[:needlen], nil
							}
						}
					}
				}
			}
			return
		})
		for scanner.Scan() {
			err := scanner.Err()
			if err != nil {
				t.onError(conn, err)
				break
			}
			id, exists := nets.GetInfoIdByConn(conn)
			buf := scanner.Bytes()
			var temp []byte = make([]byte, 0, len(buf)-PCK_MIN_SIZE)
			datas := bytes.NewBuffer(temp)
			datas.Write(buf[PCK_MIN_SIZE:])
			if exists {
				t.onMsg(conn, id, datas.Bytes())
			} else {
				if err := t.dealHandShake(conn, string(datas.Bytes())); err != nil {
					t.onError(conn, err)
				}
			}
		}
	*/
}

func (t *TcpNetWorker) Connect(id string, url string, origin string) error {

	theUrl := strings.Trim(url, "tcp://") // trim the ws header
	infos := strings.Split(theUrl, "/")   // parse the sub path

	tcpAddr, err := net.ResolveTCPAddr("tcp", infos[0])
	if err != nil {
		return err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err == nil {
		t.doHandShake(conn, origin, url, id)
		go t.h_tcpSocket(conn)
	}
	return err
}

func (t *TcpNetWorker) Send(conn net.Conn, msg []byte) error {
	datalen := len(msg)
	buf, err := binbuf.BuildBuffer(datalen + PCK_MIN_SIZE)
	if err != nil {
		return err
	}
	defer func() {
		msg = nil // dispose the send buffer
		buf.Dispose()
		buf = nil
	}()
	buf.PushInt(PCK_HEADER)
	buf.PushShort(int16(datalen))
	buf.PushBytes(msg)
	_, err2 := conn.Write(buf.Bytes())
	return err2
}

func (t *TcpNetWorker) onConn(conn net.Conn, id string) {
	// record the set from id to conn
	err := nets.AddConnInfo(id, nets.TCP, conn, t)
	if err != nil {
		t.onError(conn, err)
	} else {
		t.eventListener.OnConn(id)
	}
}

func (t *TcpNetWorker) onMsg(conn net.Conn, id string, msg []byte) {
	nets.ResetConnState(id)
	if msg[0] == 35 { // '#'
		strmsg := string(msg)
		if strmsg == "#pong" {
			fmt.Println("receive pong from.." + id) // nothing to do but ResetConnState for AutoPing
			return
		} else if strmsg == "#ping" {
			fmt.Println("re sending pong to..." + id)
			go t.Send(conn, []byte("#pong")) // return the pong pck
			return
		}
	}
	t.eventListener.OnMsg(id, msg)
}

func (t *TcpNetWorker) onClose(id string, conn net.Conn, reason error) {
	if id != "" {
		t.eventListener.OnClose(id, reason)
		nets.RemoveConnInfo(id) // remove the closed conn from local record
	}
	conn.Close()
}

func (t *TcpNetWorker) onError(conn net.Conn, err error) {
	if conn != nil {
		id, exists := nets.GetInfoIdByConn(conn)
		if exists {
			t.eventListener.OnError(id, err)
			t.onClose(id, conn, err) // close the conn with errors
		} else {
			t.eventListener.OnError("", err)
			t.onClose(id, conn, err) // close the conn with errors
		}
	} else {
		t.eventListener.OnError("", err)
	}
}

func (t *TcpNetWorker) BindEventListener(eventListener nets.INetEventListener) error {
	if t.eventListener == nil {
		t.eventListener = eventListener
		return nil
	}
	return errors.New("t net worker has binded an event listener!!")
}

func (t *TcpNetWorker) Close(id string, conn net.Conn) error {
	t.eventListener.OnClose(id, nil)
	nets.RemoveConnInfo(id) // remove the closed conn from local record
	return conn.Close()
}

func (t *TcpNetWorker) doHandShake(conn net.Conn, origin string, url string, id string) error {
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
	t.onConn(conn, id)
	return nil
}

func (t *TcpNetWorker) dealHandShake(conn net.Conn, info string) error {
	var datas map[string]string
	if err := jsoniter.Unmarshal([]byte(info), &datas); err != nil {
		return err
	}
	origin, exists := datas["Origin"]
	if !exists {
		return errors.New("handshake datas missing!")
	}
	id, b := t.eventListener.OnCheckNode(origin) // let the gonode to check if the url is legal
	if b {
		t.onConn(conn, id)
		return nil
	} else {
		return errors.New("handshake illegal!!")
	}
}
