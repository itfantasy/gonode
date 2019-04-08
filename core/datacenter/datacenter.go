package datacenter

import (
	"errors"

	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/components/etcd"
	"github.com/itfantasy/gonode/components/redis"
)

type IDCCallbacks interface {
	OnNewNode(id string)
	OnError(err error)
}

type IDataCenter interface {
	BindCallbacks(IDCCallbacks)
	RegisterAndDetect(*gen_server.NodeInfo, string, int)
	GetNodeInfo(string) (*gen_server.NodeInfo, error)
	CheckNode(string, string) bool
}

func NewDataCenter(comp interface{}) (IDataCenter, error) {
	switch comp.(type) {
	case *redis.Redis:
		return NewRedisDC(comp.(*redis.Redis)), nil
	case *etcd.Etcd:
		return NewEtcdDC(comp.(*etcd.Etcd)), nil
	}
	return nil, errors.New("illegal DC comp type! only etcd or redis ... ")
}
