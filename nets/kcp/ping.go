package kcp

import (
	"fmt"
	"sync"

	"github.com/itfantasy/gonode/utils/timer"
	"github.com/itfantasy/gonode/utils/ts"
)

/*
use a pingpong package to check if the kcpconn is alive
1. record the ts of the last recivingmsg for every conn
2. when one conn have a long enmpty time, then send a pingpong pck to check if the conn is alive
*/

type ConnState struct {
	ping bool  // if has sended a ping pck
	ts   int64 // the ts of the last recivingmsg
}

var connStates map[string]*ConnState
var stateLock sync.RWMutex

func (this *KcpNetWorker) initStates() {
	connStates = make(map[string]*ConnState)
}

func (this *KcpNetWorker) autoPing() {
	dirtyIds := make([]string, 0, 100)
	for {
		ms := ts.MS()
		timer.Sleep(3000)
		stateLock.Lock()
		for id, state := range connStates {
			if ms-state.ts > 5000 {
				if state.ping {
					dirtyIds = append(dirtyIds, id)
					fmt.Println("conn time out.." + id)
				} else if ms-state.ts > 10000 {
					fmt.Println("sending a ping to..." + id)
					this.SendAsync(id, []byte("#ping"))
					state.ping = true
				}
			}
			timer.Sleep(10)
		}
		stateLock.Unlock()
		for _, dirtyId := range dirtyIds {
			this.Close(dirtyId)
		}
		dirtyIds = dirtyIds[0:0]
	}
}

func (this *KcpNetWorker) newConnState(id string) {
	stateLock.Lock()
	defer stateLock.Unlock()

	connStates[id] = new(ConnState)
}

func (this *KcpNetWorker) disposeConnState(id string) {
	stateLock.Lock()
	defer stateLock.Unlock()

	delete(connStates, id)
}

func (this *KcpNetWorker) resetConnState(id string) {
	stateLock.Lock()
	defer stateLock.Unlock()

	state, exist := connStates[id]
	if exist {
		state.ping = false
		state.ts = ts.MS()
	}
}
