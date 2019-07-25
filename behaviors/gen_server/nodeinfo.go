package gen_server

import (
	"fmt"

	"github.com/itfantasy/gonode/utils/crypt"
	"github.com/itfantasy/gonode/utils/snowflake"
)

// the node information
type NodeInfo struct {
	Id       string
	Url      string
	Pub      bool
	BackEnds string

	LogLevel string
	LogComp  string

	RegComp string

	Sig string

	UserDatas map[string]string
}

func NewNodeInfo() *NodeInfo {
	info := new(NodeInfo)
	info.UserDatas = make(map[string]string)
	return info
}

func (info *NodeInfo) Signature() {
	info.Sig = crypt.Md5(info.Url + "?id=" + info.Id + "&sf=" + snowflake.Generate()) // call when register node info to DC, then other node use the sig to connect to this node <----> check url
}

func (info *NodeInfo) ToString() string {
	str := "NodeInfo:\r\nId:" + info.Id + "\r\nUrl:" + info.Url + "\r\nPub:" + fmt.Sprint(info.Pub) + "\r\nBackEnds:" + info.BackEnds
	return str
}
