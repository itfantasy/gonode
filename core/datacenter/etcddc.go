package datacenter

import (
	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/components/etcd"
)

type EtcdDC struct {
	coreEtcd  *etcd.Etcd
	callbacks IDCCallbacks
	info      *gen_server.NodeInfo
	channel   string
}

func NewEtcdDC(et *etcd.Etcd) *EtcdDC {
	this := new(EtcdDC)
	this.coreEtcd = et
	return this
}

func (this *EtcdDC) BindCallbacks(callbacks IDCCallbacks) {
	this.callbacks = callbacks
}
func (this *EtcdDC) RegisterAndDetect(info *gen_server.NodeInfo, channel string, msfrequency int) {
	this.info = info
	this.channel = channel
}
func (this *EtcdDC) GetNodeInfo(id string) (*gen_server.NodeInfo, error) {
	if this.info.Id == id {
		return this.info, nil
	}

	return nil, nil
}
func (this *EtcdDC) CheckNode(id string, origin string) bool {
	return false
}

func (this *EtcdDC) OnSubscribe(channel string) {

}
func (this *EtcdDC) OnSubMessage(channel string, msg string) {

}
func (this *EtcdDC) OnSubError(channel string, err error) {

}
