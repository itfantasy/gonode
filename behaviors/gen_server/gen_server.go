package gen_server

type GenServer interface {
	SelfInfo() (*NodeInfo, error) // the node self information

	Start()                 // when start
	OnDetect(string) bool   // detect a new node
	OnConn(string)          // get a new conn
	OnMsg(string, []byte)   // receive a new msg
	OnClose(string)         // a conn disconnected
	OnShell(string, string) // receive a pub/sub msg from redis
	OnRanId() string        // create a random conn id when the node is wan
}
