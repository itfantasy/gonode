package ws

import (
	"errors"
	"net/http"
	"strings"

	"net"

	"github.com/itfantasy/gonode/nets"
	"golang.org/x/net/websocket"
)

type WSNetWorker struct {
	eventListener nets.INetEventListener
}

func (this *WSNetWorker) Listen(url string) error {
	nets.InitKvvk()
	url = strings.Trim(url, "ws://") // trim the ws header
	infos := strings.Split(url, "/") // parse the sub path
	http.Handle("/"+infos[1], websocket.Handler(this.h_webSocket))
	err := http.ListenAndServe(infos[0], nil)
	return err
}

func (this *WSNetWorker) h_webSocket(conn *websocket.Conn) {
	origin := conn.RemoteAddr().String()
	urlAndId := strings.Split(origin, "#")
	if len(urlAndId) != 2 {
		err := errors.New("illegal origin data! " + origin)
		this.onError(conn, err)
		return
	}
	url := urlAndId[0]
	id := urlAndId[1]
	id, b := this.eventListener.OnCheckNode(id, url) // let the gonode to check if the url is legal
	if b {
		this.onConn(conn, id)
		var msg []byte
		for {
			err := websocket.Message.Receive(conn, &msg)
			if err != nil {
				this.onError(conn, err)
				break
			}
			this.onMsg(conn, msg)
		}
	} else {
		conn.Close() // dispose the illegel conn
	}
}

func (this *WSNetWorker) Connect(id string, url string, origin string) error {
	nets.InitKvvk()
	id, b := this.eventListener.OnCheckNode(id, url)
	if !b {
		return errors.New("handshake illegal!! " + url + "#" + id)
	}
	conn, err := websocket.Dial(url, "tcp", origin)
	if err == nil {
		go this.h_webSocket(conn)
	}
	return err
}

func (this *WSNetWorker) Send(conn net.Conn, msg []byte) error {
	defer func() {
		msg = nil // dispose the send buffer
	}()
	err := websocket.Message.Send(conn.(*websocket.Conn), msg)
	return err
}

func (this *WSNetWorker) onConn(conn *websocket.Conn, id string) {
	// record the set from id to conn
	err := nets.AddConnInfo(id, nets.WS, conn, this)
	if err != nil {
		this.onError(conn, err)
	} else {
		this.eventListener.OnConn(id)
	}
}

func (this *WSNetWorker) onMsg(conn *websocket.Conn, msg []byte) {
	id, exists := nets.GetInfoIdByConn(conn)
	if exists {
		this.eventListener.OnMsg(id, msg)
	}
}

func (this *WSNetWorker) onClose(conn *websocket.Conn) {
	id, exists := nets.GetInfoIdByConn(conn)
	if exists {
		this.eventListener.OnClose(id)
		nets.RemoveConnInfo(id) // remove the closed conn from local record
		conn.Close()
	}
}

func (this *WSNetWorker) onError(conn *websocket.Conn, err error) {
	id, exists := nets.GetInfoIdByConn(conn)
	if exists {
		this.eventListener.OnError(id, err)
		this.onClose(conn) // close the conn with errors
	}
}

func (this *WSNetWorker) BindEventListener(eventListener nets.INetEventListener) error {
	if this.eventListener == nil {
		this.eventListener = eventListener
		return nil
	}
	return errors.New("this net worker has binded an event listener!!")
}

func (this *WSNetWorker) Close(id string, conn net.Conn) error {
	this.eventListener.OnClose(id)
	nets.RemoveConnInfo(id) // remove the closed conn from local record
	return conn.Close()
}
