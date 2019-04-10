package datacenter

import (
	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/components/etcd"
	"github.com/itfantasy/gonode/utils/json"
	"github.com/itfantasy/gonode/utils/timer"
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
func (this *EtcdDC) RegisterAndDetect(info *gen_server.NodeInfo, channel string, msfrequency int) error {
	this.info = info
	this.channel = channel

	// sub the channel
	this.coreEtcd.BindSubscriber(this)
	this.coreEtcd.Subscribe(channel)

	// register self
	infoStr, err := json.Encode(this.info)
	if err != nil {
		return err
	}

	err2 := this.coreEtcd.Set(channel+"/"+this.info.Id, infoStr)
	if err2 != nil {
		return err2
	}

	// and auto detect per msfrequency
	go func() {
		for {
			timer.Sleep(msfrequency)
			ids, err := this.coreEtcd.Gets(channel + ":all")
			if err != nil {
				this.callbacks.OnError(err)
				continue
			}
			for id, _ := range ids {
				if id != this.info.Id {
					this.callbacks.OnNewNode(id)
				}
			}
		}
	}()

	return nil
}
func (this *EtcdDC) GetNodeInfo(id string) (*gen_server.NodeInfo, error) {
	if this.info.Id == id {
		return this.info, nil
	}

	return nil, nil
}
func (this *EtcdDC) CheckNode(id string, origin string) bool {
	info, err := this.GetNodeInfo(id)
	if err != nil {
		return false
	}
	url := info.Url
	if extractIPFromUrl(url) != extractIPFromUrl(origin) {
		return false
	}
	return true
}

func (this *EtcdDC) OnSubscribe(channel string) {

}
func (this *EtcdDC) OnSubMessage(channel string, msg string) {

}
func (this *EtcdDC) OnSubError(channel string, err error) {

}
