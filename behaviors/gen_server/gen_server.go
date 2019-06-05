package gen_server

type GenServer interface {
	Setup() *NodeInfo      // the node self information
	Start()                // when start
	OnConn(string)         // get a new conn
	OnMsg(string, []byte)  // receive a new msg
	OnClose(string, error) // a conn disconnected
}
