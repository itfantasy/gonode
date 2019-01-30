package gen_server

// the node information
type NodeInfo struct {
	Id         string // the index, auto create if null
	Label      string
	Url        string // the node url
	AutoDetect bool   // is auto detecting
	BackEnds   string
	Public     bool
	PubUrls    string

	RedUrl  string // the core redis url
	RedPool int    // max conn num of redis
	RedDB   int    // the db
	RedAuth string // auto info
}

func NewNodeInfo() *NodeInfo {
	return new(NodeInfo)
}
