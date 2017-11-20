package gonode

func (this *GoNode) OnConn(id string) {
	this.behavior.OnConn(id)
}

func (this *GoNode) OnMsg(id string, msg []byte) {
	this.behavior.OnMsg(id, msg)
}

func (this *GoNode) OnClose(id string) {
	this.behavior.OnClose(id)
}

func (this *GoNode) OnError(id string, err error) {
	this.logger.Error(this.sprinfLog("the node[" + id + "] occurs errors:" + err.Error()))
}
