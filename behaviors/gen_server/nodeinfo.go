package gen_server

import (
	"fmt"

	"github.com/itfantasy/gonode/utils/crypt"
	"github.com/itfantasy/gonode/utils/json"
	"github.com/itfantasy/gonode/utils/snowflake"
)

// the node information
type NodeInfo struct {
	RegDC     string
	NameSpace string
	NodeId    string
	EndPoints []string
	IsPub     bool
	BackEnds  string
	Sig       string
	UsrDatas  map[string]string
}

func NewNodeInfo() *NodeInfo {
	info := new(NodeInfo)
	info.EndPoints = make([]string, 0, 0)
	info.UsrDatas = make(map[string]string)
	return info
}

func (info *NodeInfo) Signature() {
	info.Sig = crypt.Md5("?nodeid=" + info.NodeId + "&sf=" + snowflake.Generate()) // call when register node info to DC, then other node use the sig to connect to this node <----> check url
}

func (info *NodeInfo) ToString() string {
	endPointsStr, _ := json.Marshal(info.EndPoints)
	str := " ================= - Info - ================= \r\n nodeId:" + info.NodeId + "\r\n endPoints:" + endPointsStr + "\r\n isPub:" + fmt.Sprint(info.IsPub) + "\r\n backEnds:" + info.BackEnds + "\r\n" + " ============================================ \r\n"
	return str
}
