package gen_server

// the node information
type NodeInfo struct {
	Id       string
	Url      string
	Pub      bool
	BackEnds string

	LogLevel string
	LogComp  string

	RegComp string

	UserDatas map[string]string
}

func NewNodeInfo() *NodeInfo {
	info := new(NodeInfo)
	info.UserDatas = make(map[string]string)
	return info
}
