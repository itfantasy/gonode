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
	this.initKvvk()
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
	id, exists := this.getInfoIdByConn(conn)
	//url := conn.RemoteAddr().String()
	//id, legal := this.eventListener.CheckUrlLegal(url) // let the gonode to check if the url is legal
	//if legal {
	//this.onConn(conn, id)
	var buf [4096]byte // the rev buf
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			this.onError(conn, err)
			break
		}
		if n > 0 {
			var temp []byte = make([]byte, n)
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
	//} else {
	//conn.Close() // dispose the illegel conn
	//}
}

func (this *TcpNetWorker) Connect(url string, origin string) error {
	this.initKvvk()

	url = strings.Trim(url, "tcp://") // trim the ws header
	infos := strings.Split(url, "/")  // parse the sub path

	tcpAddr, err := net.ResolveTCPAddr("tcp", infos[0])
	if err != nil {
		return err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err == nil {
		this.doHandShake(conn, origin)
		go this.h_tcpSocket(conn)
	}
	return err
}

func (this *TcpNetWorker) Send(id string, msg []byte) error {
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

func (this *TcpNetWorker) SendAsync(id string, msg []byte) {
	go this.Send(id, msg)
}

func (this *TcpNetWorker) onConn(conn net.Conn, id string) {
	// record the set from id to conn
	err := this.addConnInfo(id, conn)
	if err != nil {
		this.onError(conn, err)
	} else {
		this.eventListener.OnConn(id)
	}
}

func (this *TcpNetWorker) onMsg(conn net.Conn, id string, msg []byte) {
	//id, exists := this.getInfoIdByConn(conn)
	//if exists {
	this.eventListener.OnMsg(id, msg)
	//}
}

func (this *TcpNetWorker) onClose(conn net.Conn) {
	id, exists := this.getInfoIdByConn(conn)
	if exists {
		this.eventListener.OnClose(id)
		this.removeConnInfo(id) // remove the closed conn from local record
		conn.Close()
	}
}

func (this *TcpNetWorker) onError(conn net.Conn, err error) {
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

func (this *TcpNetWorker) BindEventListener(eventListener nets.INetEventListener) error {
	if this.eventListener == nil {
		this.eventListener = eventListener
		return nil
	}
	return errors.New("this net worker has binded an event listener!!")
}

func (this *TcpNetWorker) Close(id string) error {
	conn, exists := this.getInfoConnById(id)
	if exists {
		this.eventListener.OnClose(id)
		this.removeConnInfo(id) // remove the closed conn from local record
		return conn.Close()
	}
	return errors.New("there is not the id in local record!")
}

func (this *TcpNetWorker) doHandShake(conn net.Conn, origin string) error {
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
	url := origin
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
