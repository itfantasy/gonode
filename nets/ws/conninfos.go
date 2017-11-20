package ws

import (
	"errors"
	"sync"

	"golang.org/x/net/websocket"
)

type ConnInfos struct {
	kv   map[string]*websocket.Conn
	vk   map[*websocket.Conn]string
	LOCK sync.RWMutex
}

var connInfos *ConnInfos

func (this *WSNetWorker) initKvvk() {
	if connInfos == nil {
		connInfos = new(ConnInfos)
		connInfos.kv = make(map[string]*websocket.Conn)
		connInfos.vk = make(map[*websocket.Conn]string)
	}
}

func (this *WSNetWorker) addConnInfo(id string, conn *websocket.Conn) error {
	connInfos.LOCK.Lock()
	defer connInfos.LOCK.Unlock()

	_, ok := connInfos.kv[id]
	_, ok2 := connInfos.vk[conn]
	if ok || ok2 {
		return errors.New("a same conn info has existed!")
	}
	connInfos.kv[id] = conn
	connInfos.vk[conn] = id

	return nil
}

func (this *WSNetWorker) removeConnInfo(id string) {
	conn, ok := connInfos.kv[id]
	_, ok2 := connInfos.vk[conn]

	connInfos.LOCK.Lock()
	defer connInfos.LOCK.Unlock()

	if ok {
		delete(connInfos.kv, id)
	}
	if ok2 {
		delete(connInfos.vk, conn)
	}

}

func (this *WSNetWorker) getInfoIdByConn(conn *websocket.Conn) (string, bool) {
	connInfos.LOCK.Lock()
	defer connInfos.LOCK.Unlock()

	val, exist := connInfos.vk[conn]
	return val, exist
}

func (this *WSNetWorker) getInfoConnById(id string) (*websocket.Conn, bool) {
	connInfos.LOCK.Lock()
	defer connInfos.LOCK.Unlock()

	val, exist := connInfos.kv[id]
	return val, exist
}

func (this *WSNetWorker) GetAllConnIds() []string {
	connInfos.LOCK.Lock()
	defer connInfos.LOCK.Unlock()

	sorted_keys := make([]string, 0, len(connInfos.kv)) // set the capacity
	for k, _ := range connInfos.kv {
		sorted_keys = append(sorted_keys, k) // and this append will not create extra memory costing
	}
	//sort.Strings(sorted_keys)
	return sorted_keys
}

func (this *WSNetWorker) IsIdExists(id string) bool {
	_, exists := this.getInfoConnById(id)
	return exists
}
