package datacenter

import (
	"github.com/itfantasy/gonode/behaviors/gen_server"
)

type DataCenterPrivilege struct {
	dc IDataCenter
}

func NewDataCenterPrivilege(dc IDataCenter) *DataCenterPrivilege {
	d := new(DataCenterPrivilege)
	d.dc = dc
	return d
}

func (d *DataCenterPrivilege) GetNodeInfo(id string) (*gen_server.NodeInfo, error) {
	return d.dc.GetNodeInfo(id)
}
func (d *DataCenterPrivilege) GetNodeStatus(id string, ref interface{}) error {
	return d.dc.GetNodeStatus(id, ref)
}
