package datacenter

import (
	"strings"

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
	go this.coreEtcd.Subscribe(channel)

	// register self
	this.info.Signature()
	infoStr, err := json.Encode(this.info)
	if err != nil {
		return err
	}

	err2 := this.coreEtcd.Set(channel+"/"+this.info.Id, infoStr)
	if err2 != nil {
		return err2
	}

	// and auto detect per msfrequency
	if this.info.BackEnds != "" {
		go func() {
			for {
				timer.Sleep(msfrequency)
				ids, err := this.coreEtcd.Gets(channel)
				if err != nil {
					this.callbacks.OnDCError(err)
					continue
				}
				for id, _ := range ids {
					if id != this.info.Id {
						this.callbacks.OnNewNode(id)
					}
				}
			}
		}()
	}

	return nil
}
func (this *EtcdDC) GetNodeInfo(id string) (*gen_server.NodeInfo, error) {
	if this.info.Id == id {
		return this.info, nil
	}
	infoStr, err := this.coreEtcd.Get(this.channel + "/" + id)
	if err != nil {
		return nil, err
	}
	var info gen_server.NodeInfo
	err2 := json.Decode(infoStr, &info)
	if err2 != nil {
		return nil, err2
	}
	return &info, nil
}
func (this *EtcdDC) GetNodeSig(id string) (string, error) {
	info, err := this.GetNodeInfo(id)
	if err != nil {
		return "", err
	}
	return info.Sig, err
}
func (this *EtcdDC) CheckNode(id string, sig string) bool {
	nodeSig, err := this.GetNodeSig(id)
	if err != nil {
		return false
	}
	if nodeSig != sig {
		return false
	}
	return true
}

func (this *EtcdDC) OnSubscribe(path string) {
	if this.channel == path {

	}
}
func (this *EtcdDC) OnSubMessage(path string, msg string) {
	if strings.HasPrefix(path, this.channel) {
		this.callbacks.OnNewNode(strings.TrimPrefix(path, this.channel+"/"))
	}
}
func (this *EtcdDC) OnSubError(path string, err error) {
	if strings.HasPrefix(path, this.channel) {
		this.callbacks.OnDCError(err)
	}
}
