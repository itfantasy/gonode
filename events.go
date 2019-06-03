package gonode

import (
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
	e.node.onConn(id)
}

func (e *EventHandler) OnMsg(id string, msg []byte) {
	e.node.onMsg(id, msg)
}

func (e *EventHandler) OnClose(id string) {
	e.node.onClose(id)
}

func (e *EventHandler) OnError(id string, err error) {
	e.node.onError(id, err)
}

func (e *EventHandler) OnCheckNode(origin string) (string, bool) {
	return e.node.onCheckNode(origin)
}

func (e *EventHandler) OnNewNode(id string) {
	e.node.onNewNode(id)
}

func (e *EventHandler) OnDCError(err error) {
	e.node.onDCError(err)
}

func (e *EventHandler) OnReportError(err interface{}) {
	e.node.reportError(err)
}

func (g *GoNode) onConn(id string) {
	defer g.autoRecover()
	g.logger.Info("conn to " + id + " succeed!")
	g.behavior.OnConn(id)
}

func (g *GoNode) onMsg(id string, msg []byte) {
	defer g.autoRecover()
	g.behavior.OnMsg(id, msg)
}

func (g *GoNode) onClose(id string) {
	defer g.autoRecover()
	g.behavior.OnClose(id)
}

func (g *GoNode) onError(id string, err error) {
	defer g.autoRecover()
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
			connId := g.randomCntId()
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
func (g *GoNode) onNewNode(id string) {
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
				err2 := g.Connnect(info.Id, info.Url)
				if err2 != nil {
					g.logger.Error(err2.Error() + "[" + id + "]")
				}
			} else {
				g.logger.Error(err.Error() + "[" + id + "]")
			}
		}
	}
}

func (g *GoNode) onDCError(err error) {
	g.logger.Error(err.Error())
}

func (g *GoNode) onDispose() {

}
