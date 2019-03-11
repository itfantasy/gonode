package gonode

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
