package kcp

import (
	"errors"
	"sync"

	"net"
)

type ConnInfos struct {
	kv   map[string]net.Conn
	vk   map[net.Conn]string
	LOCK sync.RWMutex
}

var connInfos *ConnInfos

func (this *KcpNetWorker) initKvvk() {
	if connInfos == nil {
		connInfos = new(ConnInfos)
		connInfos.kv = make(map[string]net.Conn)
		connInfos.vk = make(map[net.Conn]string)
	}
}

func (this *KcpNetWorker) addConnInfo(id string, conn net.Conn) error {
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

func (this *KcpNetWorker) removeConnInfo(id string) {
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

func (this *KcpNetWorker) getInfoIdByConn(conn net.Conn) (string, bool) {
	connInfos.LOCK.Lock()
	defer connInfos.LOCK.Unlock()

	val, exist := connInfos.vk[conn]
	return val, exist
}

func (this *KcpNetWorker) getInfoConnById(id string) (net.Conn, bool) {
	connInfos.LOCK.Lock()
	defer connInfos.LOCK.Unlock()

	val, exist := connInfos.kv[id]
	return val, exist
}

func (this *KcpNetWorker) GetAllConnIds() []string {
	connInfos.LOCK.Lock()
	defer connInfos.LOCK.Unlock()

	sorted_keys := make([]string, 0, len(connInfos.kv)) // set the capacity
	for k, _ := range connInfos.kv {
		sorted_keys = append(sorted_keys, k) // and this append will not create extra memory costing
	}
	//sort.Strings(sorted_keys)
	return sorted_keys
}

func (this *KcpNetWorker) IsIdExists(id string) bool {
	_, exists := this.getInfoConnById(id)
	return exists
}
