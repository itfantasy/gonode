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

func (w *WSNetWorker) Listen(url string) error {

	url = strings.Trim(url, "ws://") // trim the ws header
	infos := strings.Split(url, "/") // parse the sub path
	http.Handle("/"+infos[1], websocket.Handler(w.h_webSocket))
	err := http.ListenAndServe(infos[0], nil)
	return err
}

func (w *WSNetWorker) h_webSocket(conn *websocket.Conn) {
	remote := conn.RemoteAddr().String()
	id, b := w.eventListener.OnCheckNode(remote) // let the gonode to check if the url is legal
	if b {
		w.onConn(conn, id)
		var msg []byte
		for {
			err := websocket.Message.Receive(conn, &msg)
			if err != nil {
				w.onError(conn, err)
				break
			}
			w.onMsg(conn, msg)
		}
	} else {
		conn.Close() // dispose the illegel conn
	}
}

func (w *WSNetWorker) Connect(id string, url string, origin string) error {
	conn, err := websocket.Dial(url, "tcp", origin)
	if err == nil {
		w.onConn(conn, id)
		var msg []byte
		for {
			err := websocket.Message.Receive(conn, &msg)
			if err != nil {
				w.onError(conn, err)
				break
			}
			w.onMsg(conn, msg)
		}
	}
	return err
}

func (w *WSNetWorker) Send(conn net.Conn, msg []byte) error {
	defer func() {
		msg = nil // dispose the send buffer
	}()
	err := websocket.Message.Send(conn.(*websocket.Conn), msg)
	return err
}

func (w *WSNetWorker) onConn(conn *websocket.Conn, id string) {
	// record the set from id to conn
	err := nets.AddConnInfo(id, nets.WS, conn, w)
	if err != nil {
		w.onError(conn, err)
	} else {
		w.eventListener.OnConn(id)
	}
}

func (w *WSNetWorker) onMsg(conn *websocket.Conn, msg []byte) {
	id, exists := nets.GetInfoIdByConn(conn)
	if exists {
		w.eventListener.OnMsg(id, msg)
	}
}

func (w *WSNetWorker) onClose(conn *websocket.Conn) {
	id, exists := nets.GetInfoIdByConn(conn)
	if exists {
		w.eventListener.OnClose(id)
		nets.RemoveConnInfo(id) // remove the closed conn from local record
		conn.Close()
	}
}

func (w *WSNetWorker) onError(conn *websocket.Conn, err error) {
	id, exists := nets.GetInfoIdByConn(conn)
	if exists {
		w.eventListener.OnError(id, err)
		w.onClose(conn) // close the conn with errors
	}
}

func (w *WSNetWorker) BindEventListener(eventListener nets.INetEventListener) error {
	if w.eventListener == nil {
		w.eventListener = eventListener
		return nil
	}
	return errors.New("w net worker has binded an event listener!!")
}

func (w *WSNetWorker) Close(id string, conn net.Conn) error {
	w.eventListener.OnClose(id)
	nets.RemoveConnInfo(id) // remove the closed conn from local record
	return conn.Close()
}
