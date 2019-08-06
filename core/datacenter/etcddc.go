package datacenter

import (
	"strings"

	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/components/etcd"
	"github.com/itfantasy/gonode/nets"
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
	infoStr, err := json.Encode(e.info)
	if err != nil {
		return err
	}

	err2 := e.coreEtcd.Set(channel+"/"+e.info.Id, infoStr)
	if err2 != nil {
		return err2
	}

	// and auto detect per msfrequency
	if e.info.BackEnds != "" {
		go func() {
			for {
				timer.Sleep(msfrequency)
				ids, err := e.coreEtcd.Gets(channel)
				if err != nil {
					e.callbacks.OnDCError(err)
					continue
				}
				for idstr, _ := range ids {
					id := strings.TrimPrefix(idstr, channel+"/")
					if id != e.info.Id {
						connErr := e.callbacks.OnNewNode(id)
						if e.info.Id == "supervisor" {
							if connErr != nil {
								if checkOutOfDate(id) {
									if _, err := e.coreEtcd.Delete(idstr); err != nil {
										e.callbacks.OnDCError(err)
									} else if _, err := e.coreEtcd.Delete(channel + "-Status/" + id); err != nil {
										e.callbacks.OnDCError(err)
									} else {
										clearOutOfDate(id)
										e.callbacks.OnNodeDestruct(id)
									}
								}
							} else {
								clearOutOfDate(id)
							}
						}
					}
				}
				if e.info.Id != "supervisor" {
					e.updateNodeStatus(nets.AllSvcConnIds())
				}
			}
		}()
	}

	return nil
}
func (e *EtcdDC) GetNodeInfo(id string) (*gen_server.NodeInfo, error) {
	if e.info.Id == id {
		return e.info, nil
	}
	infoStr, err := e.coreEtcd.Get(e.channel + "/" + id)
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
func (e *EtcdDC) updateNodeStatus(conns []string) error {
	status, err := json.Encode(conns)
	if err != nil {
		return err
	}
	err2 := e.coreEtcd.Set(e.channel+"-Status/"+e.info.Id, status)
	if err2 != nil {
		return err2
	}
	return nil
}
func (e *EtcdDC) GetNodeStatus(id string) ([]string, error) {
	status, err := e.coreEtcd.Get(e.channel + "-Status/" + id)
	if err != nil {
		return nil, err
	}
	conns := make([]string, 0, 0)
	err2 := json.Decode(status, conns)
	if err2 != nil {
		return nil, err2
	}
	return conns, nil
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
