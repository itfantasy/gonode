package supervisor

type Supervisor interface {
	Setup() *SuperInfo
	Start()                // when start
	OnConn(string)         // get a new conn
	OnConnFailed(string)   // when try to conn failed
	OnMsg(string, []byte)  // receive a new msg
	OnClose(string, error) // a conn disconnected
	OnUnregister(string)   // a node has been unregister (clean from the datacenter)
}
