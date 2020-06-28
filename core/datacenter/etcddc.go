package datacenter

import (
	"strings"

	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/components"
	"github.com/itfantasy/gonode/utils/json"
	"github.com/itfantasy/gonode/utils/timer"
)

type EtcdDC struct {
	coreEtcd  *components.Etcd
	callbacks IDCCallbacks
	info      *gen_server.NodeInfo
	channel   string
}

func NewEtcdDC(et *components.Etcd) *EtcdDC {
	e := new(EtcdDC)
	e.coreEtcd = et
	return e
}

func (e *EtcdDC) BindCallbacks(callbacks IDCCallbacks) {
	e.callbacks = callbacks
}
func (e *EtcdDC) RegisterAndDetect(info *gen_server.NodeInfo, channel string, msfrequency int) error {
	e.info = info
	e.channel = channel

	// sub the channel
	e.coreEtcd.BindSubscriber(e)
	go e.coreEtcd.Subscribe(channel)

	// register self
	e.info.Signature()
	infoStr, err := json.Marshal(e.info)
	if err != nil {
		return err
	}

	if len(e.info.EndPoints) > 0 {
		if err := e.coreEtcd.Set(channel+"/"+e.info.NodeId, infoStr); err != nil {
			return err
		}
	}

	// and auto detect per msfrequency

	go func() {
		for {
			timer.Sleep(msfrequency)
			if e.info.BackEnds != "" {
				ids, err := e.coreEtcd.Gets(channel)
				if err != nil {
					e.callbacks.OnDCError(err)
					continue
				}
				for idstr, _ := range ids {
					id := strings.TrimPrefix(idstr, channel+"/")
					if id != e.info.NodeId {
						connErr := e.callbacks.OnNewNode(id)
						if e.info.NodeId == "supervisor" {
							if connErr != nil {
								if checkOutOfDate(id) {
									if _, err := e.coreEtcd.Delete(idstr); err != nil {
										e.callbacks.OnDCError(err)
									} else if _, err := e.coreEtcd.Delete("--status/" + channel + "/" + id); err != nil {
										e.callbacks.OnDCError(err)
									} else {
										clearOutOfDate(id)
										e.callbacks.OnUnregister(id)
									}
								}
							} else {
								clearOutOfDate(id)
							}
						}
					}
				}
			}
			if e.info.NodeId != "supervisor" {
				e.updateNodeStatus()
			}
		}
	}()

	return nil
}
func (e *EtcdDC) GetNodeInfo(id string) (*gen_server.NodeInfo, error) {
	if e.info.NodeId == id {
		return e.info, nil
	}
	infoStr, err := e.coreEtcd.Get(e.channel + "/" + id)
	if err != nil {
		return nil, err
	}
	var info gen_server.NodeInfo
	err2 := json.Unmarshal(infoStr, &info)
	if err2 != nil {
		return nil, err2
	}
	return &info, nil
}
func (e *EtcdDC) GetNodeSig(id string) (string, error) {
	info, err := e.GetNodeInfo(id)
	if err != nil {
		return "", err
	}
	return info.Sig, err
}
func (e *EtcdDC) CheckNode(id string, sig string) bool {
	if id == "" {
		return false
	}
	nodeSig, err := e.GetNodeSig(id)
	if err != nil {
		return false
	}
	if nodeSig != sig {
		return false
	}
	return true
}
func (e *EtcdDC) updateNodeStatus() error {
	status, err := json.Marshal(e.callbacks.OnUpdateNodeStatus())
	if err != nil {
		return err
	}
	err2 := e.coreEtcd.Set("--status/"+e.channel+"/"+e.info.NodeId, status)
	if err2 != nil {
		return err2
	}
	return nil
}
func (e *EtcdDC) GetNodeStatus(id string, ref interface{}) error {
	status, err := e.coreEtcd.Get("--status/" + e.channel + "/" + id)
	if err != nil {
		return err
	}
	err2 := json.Unmarshal(status, ref)
	if err2 != nil {
		return err2
	}
	return nil
}
func (e *EtcdDC) OnSubscribe(path string) {
	if e.channel == path {

	}
}
func (e *EtcdDC) OnSubMessage(path string, msg string) {
	if strings.HasPrefix(path, e.channel) {
		e.callbacks.OnNewNode(strings.TrimPrefix(path, e.channel+"/"))
	}
}
func (e *EtcdDC) OnSubError(path string, err error) {
	if strings.HasPrefix(path, e.channel) {
		e.callbacks.OnDCError(err)
	}
}
