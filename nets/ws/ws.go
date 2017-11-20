package ws

import (
	"errors"
	"net/http"
	"strings"

	"github.com/itfantasy/gonode/nets"
	"golang.org/x/net/websocket"
)

type WSNetWorker struct {
	eventListener nets.INetEventListener
}

func (this *WSNetWorker) Listen(url string) error {
	this.initKvvk()
	url = strings.Trim(url, "ws://") // trim the ws header
	infos := strings.Split(url, "/") // parse the sub path
	http.Handle("/"+infos[1], websocket.Handler(this.h_webSocket))
	err := http.ListenAndServe(infos[0], nil)
	return err
}

func (this *WSNetWorker) h_webSocket(conn *websocket.Conn) {
	url := conn.RemoteAddr().String()
	id, legal := this.eventListener.CheckUrlLegal(url) // let the gonode to check if the url is legal
	if legal {
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

func (this *WSNetWorker) Connect(url string, origin string) error {
	this.initKvvk()
	conn, err := websocket.Dial(url, "tcp", origin)
	if err == nil {
		go this.h_webSocket(conn)
	}
	return err
}

func (this *WSNetWorker) Send(id string, msg []byte) error {
	defer func() {
		msg = nil // dispose the send buffer
	}()
	conn, exist := this.getInfoConnById(id)
	if exist {
		err := websocket.Message.Send(conn, msg)
		return err
	} else {
		return errors.New("there is not the conn for this id in local record!")
	}
}

func (this *WSNetWorker) SendAsync(id string, msg []byte) {
	go this.Send(id, msg)
}

func (this *WSNetWorker) onConn(conn *websocket.Conn, id string) {
	// record the set from id to conn
	err := this.addConnInfo(id, conn)
	if err != nil {
		this.onError(conn, err)
	} else {
		this.eventListener.OnConn(id)
	}
}

func (this *WSNetWorker) onMsg(conn *websocket.Conn, msg []byte) {
	id, exists := this.getInfoIdByConn(conn)
	if exists {
		this.eventListener.OnMsg(id, msg)
	}
}

func (this *WSNetWorker) onClose(conn *websocket.Conn) {
	id, exists := this.getInfoIdByConn(conn)
	if exists {
		this.eventListener.OnClose(id)
		this.removeConnInfo(id) // remove the closed conn from local record
	}
}

func (this *WSNetWorker) onError(conn *websocket.Conn, err error) {
	id, exists := this.getInfoIdByConn(conn)
	if exists {
		this.eventListener.OnError(id, err)
		this.onClose(conn) // close the conn with errors
	}
}

func (this *WSNetWorker) BindEventListener(eventListener nets.INetEventListener) {
	this.eventListener = eventListener
}
