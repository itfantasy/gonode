package supervisor

type Supervisor interface {
	Setup() *SuperInfo
	Start()                // when start
	OnConn(string)         // get a new conn
	OnConnFailed(string)   // when try to conn failed
	OnMsg(string, []byte)  // receive a new msg
	OnLog(string)          // get a log from one node
	OnMonitor(string)      // get a monitor info from on node
	OnClose(string, error) // a conn disconnected
	OnUnregister(string)   // a node has been unregister (clean from the datacenter)
}
