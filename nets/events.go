package nets

type INetEventListener interface {
	OnConn(string)
	OnMsg(string, []byte)
	OnClose(string, error)
	OnError(string, error)
	OnCheckNode(string) (string, bool)
}
