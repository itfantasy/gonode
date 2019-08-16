package gonode

import (
	"github.com/itfantasy/gonode/core/erl"
	"github.com/itfantasy/gonode/nets"
)

type EventHandler struct {
	node *GoNode
}

func newEventHandler(node *GoNode) *EventHandler {
	e := new(EventHandler)
	e.node = node
	return e
}

func (e *EventHandler) OnConn(id string) {
	defer erl.AutoRecover(e)
	e.node.onConn(id)
}

func (e *EventHandler) OnMsg(id string, msg []byte) {
	defer erl.AutoRecover(e)
	e.node.onMsg(id, msg)
}

func (e *EventHandler) OnClose(id string, reason error) {
	defer erl.AutoRecover(e)
	e.node.onClose(id, reason)
}

func (e *EventHandler) OnError(id string, err error) {
	e.node.onError(id, err)
}

func (e *EventHandler) OnCheckNode(origin string) (string, bool) {
	return e.node.onCheckNode(origin)
}

func (e *EventHandler) OnNewNode(id string) error {
	return e.node.onNewNode(id)
}

func (e *EventHandler) OnDCError(err error) {
	e.node.onDCError(err)
}

func (e *EventHandler) OnUnregister(id string) {
	e.node.onUnregister(id)
}
func (e *EventHandler) OnUpdateNodeStatus() interface{} {
	return nets.AllSvcIds()
}
func (e *EventHandler) OnDigestError(err interface{}) {
	e.node.reportError(err)
}

func (g *GoNode) onConn(id string) {
	g.logger.Info("conn to " + id + " succeed!")
	g.behavior.OnConn(id)
}

func (g *GoNode) onMsg(id string, msg []byte) {
	g.behavior.OnMsg(id, msg)
}

func (g *GoNode) onClose(id string, reason error) {
	g.behavior.OnClose(id, reason)
}

func (g *GoNode) onError(id string, err error) {
	g.logger.Error("the node[" + id + "] occurs errors:" + err.Error())
}

func (g *GoNode) onCheckNode(origin string) (string, bool) {
	b := false
	id, url, sig, err := nets.ParserOriginInfo(origin)
	if err == nil {
		b = g.dc.CheckNode(id, sig)
	}
	if !b {
		if !g.info.Pub {
			g.logger.Info("not a inside node! give up the conn:" + origin)
			return "", false
		} else {
			connId := nets.RanCntId()
			return connId, true
		}
	} else {
		exist := nets.IsIdExists(id)
		if exist {
			g.logger.Info("there is a same id in local record:" + url + "#" + id)
			return "", false
		}
		return id, true
	}
}

// when a new node is found
func (g *GoNode) onNewNode(id string) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	exist := nets.IsIdExists(id)
	if !exist {
		// check the local node is interested in the new node
		if g.checkTargetId(id) {
			g.logger.Info("a new node has been found!", id)
			// find the node url by the id
			info, err := g.dc.GetNodeInfo(id)
			if err == nil {
				if connErr := g.Connnect(info.Id, info.Url); connErr != nil {
					g.logger.Error(connErr.Error() + "[" + id + "]")
					if g.info.Id == SUPERVISOR && g.super != nil {
						g.super.OnConnFailed(id)
					}
					return connErr
				}
			} else {
				g.logger.Error(err.Error() + "[" + id + "]")
			}
		}
	}
	return nil
}

func (g *GoNode) onDCError(err error) {
	g.logger.Error(err.Error())
}

func (g *GoNode) onUnregister(id string) {
	if g.info.Id == SUPERVISOR && g.super != nil {
		g.super.OnUnregister(id)
	}
}

func (g *GoNode) onDispose() {

}
