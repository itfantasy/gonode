package datacenter

import (
	"errors"
	"strings"

	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/components"
)

type IDCCallbacks interface {
	OnNewNode(id string) error
	OnDCError(err error)
	OnUnregister(id string)
	OnUpdateNodeStatus() interface{}
}

type IDataCenter interface {
	BindCallbacks(IDCCallbacks)
	RegisterAndDetect(*gen_server.NodeInfo, string, int) error
	GetNodeInfo(string) (*gen_server.NodeInfo, error)
	GetNodeSig(string) (string, error)
	CheckNode(string, string) bool
	GetNodeStatus(string, interface{}) error
}

func NewDataCenter(regcomp string) (IDataCenter, error) {
	comp, err := components.NewComponent(regcomp)
	if err != nil {
		return nil, err
	}
	switch comp.(type) {
	case *components.Redis:
		return NewRedisDC(comp.(*components.Redis)), nil
	case *components.Etcd:
		return NewEtcdDC(comp.(*components.Etcd)), nil
	}
	return nil, errors.New("illegal DC comp type! only etcd or redis ... ")
}

func extractIPFromUrl(url string) string {
	infos := strings.Split(url, "://")
	if len(infos) != 2 {
		return ""
	}
	ipAndPort := strings.Split(infos[1], ":")
	if len(ipAndPort) != 2 {
		return ""
	}
	return ipAndPort[0]
}
