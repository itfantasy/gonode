package datacenter

import (
	"strings"

	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/components/redis"
	"github.com/itfantasy/gonode/utils/json"
	"github.com/itfantasy/gonode/utils/timer"
)

const (
	GONODE_NEW_NODE string = "itfantasy.gonode.newnode"
)

type RedisDC struct {
	coreRedis *redis.Redis
	callbacks IDCCallbacks
	info      *gen_server.NodeInfo
	channel   string
}

func NewRedisDC(red *redis.Redis) *RedisDC {
	r := new(RedisDC)
	r.coreRedis = red
	return r
}

func (r *RedisDC) BindCallbacks(callbacks IDCCallbacks) {
	r.callbacks = callbacks
}
func (r *RedisDC) RegisterAndDetect(info *gen_server.NodeInfo, channel string, msfrequency int) error {
	r.info = info
	r.channel = channel

	// sub the channel
	r.coreRedis.BindSubscriber(r)
	go r.coreRedis.Subscribe(channel)

	// register self
	r.info.Signature()
	infoStr, err := json.Encode(r.info)
	if err != nil {
		return err
	}
	_, err2 := r.coreRedis.Set(channel+":infos:"+r.info.Id, infoStr)
	if err2 != nil {
		return err2
	}
	_, err3 := r.coreRedis.SAdd(channel+":all", r.info.Id)
	if err3 != nil {
		return err3
	}
	r.coreRedis.Publish(channel, GONODE_NEW_NODE+"#"+r.info.Id)

	// and auto detect per msfrequency
	if r.info.BackEnds != "" {
		go func() {
			for {
				timer.Sleep(msfrequency)
				ids, err := r.coreRedis.SMembers(channel + ":all")
				if err != nil {
					r.callbacks.OnDCError(err)
					continue
				}
				for _, id := range ids {
					if id != r.info.Id {
						r.callbacks.OnNewNode(id)
					}
				}
			}
		}()
	}

	return nil
}
func (r *RedisDC) GetNodeInfo(id string) (*gen_server.NodeInfo, error) {
	if r.info.Id == id {
		return r.info, nil
	}
	infoStr, err := r.coreRedis.Get(r.channel + ":infos:" + id)
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
func (r *RedisDC) GetNodeSig(id string) (string, error) {
	info, err := r.GetNodeInfo(id)
	if err != nil {
		return "", err
	}
	return info.Sig, err
}
func (r *RedisDC) CheckNode(id string, sig string) bool {
	if id == "" {
		return false
	}
	nodeSig, err := r.GetNodeSig(id)
	if err != nil {
		return false
	}
	if nodeSig != sig {
		return false
	}
	return true
}

func (r *RedisDC) OnSubscribe(channel string) {
	if r.channel == channel {

	}
}
func (r *RedisDC) OnSubMessage(channel string, msg string) {
	if r.channel == channel {
		infos := strings.Split(msg, "#")
		if len(infos) == 2 {
			if infos[0] == GONODE_NEW_NODE {
				r.callbacks.OnNewNode(infos[1])
			}
		}
	}
}
func (r *RedisDC) OnSubError(channel string, err error) {
	if r.channel == channel {
		r.callbacks.OnDCError(err)
	}
}
func (e *RedisDC) ApplyDestruction(id string) bool {
	// TODO: apply for destruction of one node (only supervisor)
	// maybe we should consider the appling count or time from the first apply

	return true
}
