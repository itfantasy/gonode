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

	go func() {
		for {
			timer.Sleep(msfrequency)
			if r.info.BackEnds != "" {
				ids, err := r.coreRedis.SMembers(channel + ":all")
				if err != nil {
					r.callbacks.OnDCError(err)
					continue
				}
				for _, id := range ids {
					if id != r.info.Id {
						connErr := r.callbacks.OnNewNode(id)
						if r.info.Id == "supervisor" {
							if connErr != nil {
								if checkOutOfDate(id) {
									if _, err := r.coreRedis.Delete(r.channel + ":infos:" + id); err != nil {
										r.callbacks.OnDCError(err)
									} else if _, err := r.coreRedis.Delete(r.channel + ":status:" + id); err != nil {
										r.callbacks.OnDCError(err)
									} else if _, err := r.coreRedis.SRem(r.channel+":all", id); err != nil {
										r.callbacks.OnDCError(err)
									} else {
										clearOutOfDate(id)
										r.callbacks.OnUnregister(id)
									}
								}
							} else {
								clearOutOfDate(id)
							}
						}
					}
				}
			}
			if r.info.Id != "supervisor" {
				r.updateNodeStatus()
			}
		}
	}()

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
func (r *RedisDC) updateNodeStatus() error {
	status, err := json.Encode(r.callbacks.OnUpdateNodeStatus())
	if err != nil {
		return err
	}
	_, err2 := r.coreRedis.Set(r.channel+":status:"+r.info.Id, status)
	if err2 != nil {
		return err2
	}
	return nil
}
func (r *RedisDC) GetNodeStatus(id string, ref interface{}) error {
	status, err := r.coreRedis.Get(r.channel + ":status:" + id)
	if err != nil {
		return err
	}
	err2 := json.Decode(status, ref)
	if err2 != nil {
		return err2
	}
	return nil
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
