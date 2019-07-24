package supervisor

type Supervisor interface {
	Setup() *SuperInfo
	Start()                   // when start
	OnConn(string)            // get a new conn
	OnLog(string, string)     // get a log from one node
	OnMonitor(string, string) // get a monitor info from on node
	OnClose(string, error)    // a conn disconnected
}
