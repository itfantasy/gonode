package gonode

import (
	"github.com/itfantasy/gonode/nets"
)

func (this *GoNode) OnConn(id string) {
	defer this.autoRecover()
	this.logger.Info("conn to " + id + " succeed!")
	this.behavior.OnConn(id)
}

func (this *GoNode) OnMsg(id string, msg []byte) {
	defer this.autoRecover()
	this.behavior.OnMsg(id, msg)
}

func (this *GoNode) OnClose(id string) {
	defer this.autoRecover()
	this.behavior.OnClose(id)
}

func (this *GoNode) OnError(id string, err error) {
	defer this.autoRecover()
	this.logger.Error("the node[" + id + "] occurs errors:" + err.Error())
}

func (this *GoNode) OnCheckNode(origin string) (string, bool) {
	b := false
	id, url, sig, err := nets.ParserOriginInfo(origin)
	if err == nil {
		b = this.dc.CheckNode(id, sig)
	}
	if !b {
		if !this.info.Pub {
			this.logger.Info("not a inside node! give up the conn:" + origin)
			return "", false
		} else {
			connId := this.randomCntId()
			return connId, true
		}
	} else {
		exist := nets.IsIdExists(id)
		if exist {
			this.logger.Info("there is a same id in local record:" + url + "#" + id)
			return "", false
		}
		return id, true
	}
}

// when a new node is found
func (this *GoNode) OnNewNode(id string) {
	this.lock.Lock()
	defer this.lock.Unlock()

	exist := nets.IsIdExists(id)
	if !exist {
		// check the local node is interested in the new node
		if this.checkTargetId(id) {
			this.logger.Info("a new node has been found!", id)
			// find the node url by the id
			info, err := this.dc.GetNodeInfo(id)
			if err == nil {
				err2 := this.Connnect(info.Id, info.Url)
				if err2 != nil {
					this.logger.Error(err2.Error() + "[" + id + "]")
				}
			} else {
				this.logger.Error(err.Error() + "[" + id + "]")
			}
		}
	}
}

func (this *GoNode) OnDCError(err error) {
	this.logger.Error(err.Error())
}
