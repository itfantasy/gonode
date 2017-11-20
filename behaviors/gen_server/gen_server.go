package gen_server

type GenServer interface {
	SelfNodeInfo() (*NodeInfo, error) // the node self information
	IsInterestedIn(string) bool       // is interested in the id

	Start()                 // when start
	Update()                // timer update
	OnConn(string)          // get a new conn
	OnMsg(string, []byte)   // receive a new msg
	OnClose(string)         // a conn disconnected
	OnShell(string, string) // receive a pub/sub msg from redis
	OnReload(string) error  // reload the bll
	CreateConnId() string   // create a random conn id when the node is wan
}
