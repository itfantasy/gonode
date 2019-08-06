package nets

import (
	"errors"
	"net"
	"strings"
	"sync"

	"github.com/itfantasy/gonode/utils/snowflake"
)

type ConnInfos struct {
	kv   map[string]*connItemInfo
	vk   map[net.Conn]string
	LOCK sync.RWMutex
}

type connItemInfo struct {
	id        string
	proto     string
	conn      net.Conn
	netWorker INetWorker
}

func newConnItemInfo(id string, proto string, conn net.Conn, netWorker INetWorker) *connItemInfo {
	this := new(connItemInfo)
	this.id = id
	this.proto = proto
	this.conn = conn
	this.netWorker = netWorker
	return this
}

var connInfos *ConnInfos

func init() {
	connInfos = new(ConnInfos)
	connInfos.kv = make(map[string]*connItemInfo)
	connInfos.vk = make(map[net.Conn]string)

	connStates = make(map[string]*ConnState)
}

func AddConnInfo(id string, proto string, conn net.Conn, netWorker INetWorker) error {
	connInfos.LOCK.Lock()
	defer connInfos.LOCK.Unlock()

	_, ok := connInfos.kv[id]
	_, ok2 := connInfos.vk[conn]
	if ok || ok2 {
		return errors.New("a same conn info has existed!")
	}
	connInfos.kv[id] = newConnItemInfo(id, proto, conn, netWorker)
	connInfos.vk[conn] = id

	if proto == KCP || proto == TCP {
		newConnState(id, conn)
	}
	return nil
}

func RemoveConnInfo(id string) {
	info, ok := connInfos.kv[id]
	_, ok2 := connInfos.vk[info.conn]

	connInfos.LOCK.Lock()
	defer connInfos.LOCK.Unlock()

	if ok {
		disposeConnState(id)
		delete(connInfos.kv, id)
	}
	if ok2 {
		delete(connInfos.vk, info.conn)
	}
}

func GetInfoIdByConn(conn net.Conn) (string, bool) {
	connInfos.LOCK.Lock()
	defer connInfos.LOCK.Unlock()

	val, exist := connInfos.vk[conn]
	return val, exist
}

func GetInfoConnById(id string) (net.Conn, string, INetWorker, bool) {
	connInfos.LOCK.Lock()
	defer connInfos.LOCK.Unlock()

	val, exist := connInfos.kv[id]
	if exist {
		return val.conn, val.proto, val.netWorker, exist
	} else {
		return nil, "", nil, false
	}
}

func AllConnIds() []string {
	connInfos.LOCK.Lock()
	defer connInfos.LOCK.Unlock()

	keys := make([]string, 0, len(connInfos.kv)) // set the capacity
	for k, _ := range connInfos.kv {
		keys = append(keys, k) // and this append will not create extra memory costing
	}
	return keys
}

func AllSvcConnIds() []string {
	connInfos.LOCK.Lock()
	defer connInfos.LOCK.Unlock()

	keys := make([]string, 0, len(connInfos.kv)) // set the capacity
	for k, _ := range connInfos.kv {
		if !IsCntConnId(k) {
			keys = append(keys, k) // and this append will not create extra memory costing
		}
	}
	return keys
}

func AllCntConnIds() []string {
	connInfos.LOCK.Lock()
	defer connInfos.LOCK.Unlock()

	keys := make([]string, 0, len(connInfos.kv)) // set the capacity
	for k, _ := range connInfos.kv {
		if IsCntConnId(k) {
			keys = append(keys, k) // and this append will not create extra memory costing
		}
	}
	return keys
}

func IsIdExists(id string) bool {
	connInfos.LOCK.Lock()
	defer connInfos.LOCK.Unlock()

	_, exist := connInfos.kv[id]
	return exist
}

func RanCntConnId() string {
	return "cnt-" + snowflake.Generate()
}

func Label(id string) string {
	return strings.Split(id, "-")[0]
}

func IsCntConnId(id string) bool {
	return Label(id) == "cnt"
}
