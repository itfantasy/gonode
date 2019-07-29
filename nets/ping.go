package nets

import (
	"fmt"
	"sync"

	"net"

	"github.com/itfantasy/gonode/utils/timer"
	"github.com/itfantasy/gonode/utils/ts"
)

/*
use a pingpong package to check if the kcpconn is alive
1. record the ts of the last recivingmsg for every conn
2. when one conn have a long enmpty time, then send a pingpong pck to check if the conn is alive
*/

type ConnState struct {
	id   string
	conn net.Conn
	ping bool  // if has sended a ping pck
	ts   int64 // the ts of the last recivingmsg
}

var connStates map[string]*ConnState
var stateLock sync.RWMutex

func AutoPing(netWorker INetWorker) {
	dirtyStates := make([]*ConnState, 0, 100)
	for {
		ms := ts.MS()
		timer.Sleep(2000)
		stateLock.Lock()
		for id, state := range connStates {
			if state.ping {
				dirtyStates = append(dirtyStates, state)
				fmt.Println("pingpong time out.." + id)
			} else if ms-state.ts > 2000 {
				//fmt.Println("sending a ping to..." + id)
				go netWorker.Send(state.conn, []byte("#ping"))
				state.ping = true
			}
			timer.Sleep(10)
		}
		stateLock.Unlock()
		for _, state := range dirtyStates {
			netWorker.Close(state.id, state.conn)
		}
		dirtyStates = dirtyStates[0:0]
	}
}

func newConnState(id string, conn net.Conn) {
	stateLock.Lock()
	defer stateLock.Unlock()

	state := new(ConnState)
	state.id = id
	state.conn = conn
	connStates[state.id] = state
}

func disposeConnState(id string) {
	stateLock.Lock()
	defer stateLock.Unlock()

	_, exist := connStates[id]
	if exist {
		delete(connStates, id)
	}
}

func ResetConnState(id string, netWorker INetWorker, msg []byte) bool {
	stateLock.Lock()
	defer stateLock.Unlock()

	state, exist := connStates[id]
	if exist {
		state.ping = false
		state.ts = ts.MS()
		if msg[0] == 35 { // '#'
			strmsg := string(msg)
			if strmsg == "#pong" {
				//fmt.Println("receive pong from.." + id) // nothing to do but ResetConnState for AutoPing
				return true
			} else if strmsg == "#ping" {
				//fmt.Println("re sending pong to..." + id)
				go netWorker.Send(state.conn, []byte("#pong")) // return the pong pck
				return true
			}
		}
	}
	return false
}
