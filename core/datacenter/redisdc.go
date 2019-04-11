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
	this := new(RedisDC)
	this.coreRedis = red
	return this
}

func (this *RedisDC) BindCallbacks(callbacks IDCCallbacks) {
	this.callbacks = callbacks
}
func (this *RedisDC) RegisterAndDetect(info *gen_server.NodeInfo, channel string, msfrequency int) error {
	this.info = info
	this.channel = channel

	// sub the channel
	this.coreRedis.BindSubscriber(this)
	go this.coreRedis.Subscribe(channel)

	// register self
	infoStr, err := json.Encode(this.info)
	if err != nil {
		return err
	}
	_, err2 := this.coreRedis.Set(channel+":infos:"+this.info.Id, infoStr)
	if err2 != nil {
		return err2
	}
	_, err3 := this.coreRedis.SAdd(channel+":all", this.info.Id)
	if err3 != nil {
		return err3
	}
	this.coreRedis.Publish(channel, GONODE_NEW_NODE+"#"+this.info.Id)

	// and auto detect per msfrequency
	if this.info.BackEnds != "" {
		go func() {
			for {
				timer.Sleep(msfrequency)
				ids, err := this.coreRedis.SMembers(channel + ":all")
				if err != nil {
					this.callbacks.OnDCError(err)
					continue
				}
				for _, id := range ids {
					if id != this.info.Id {
						this.callbacks.OnNewNode(id)
					}
				}
			}
		}()
	}

	return nil
}
func (this *RedisDC) GetNodeInfo(id string) (*gen_server.NodeInfo, error) {
	if this.info.Id == id {
		return this.info, nil
	}

	return nil, nil
}
func (this *RedisDC) CheckNode(id string, origin string) bool {
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

func (this *RedisDC) OnSubscribe(channel string) {
	if this.channel == channel {

	}
}
func (this *RedisDC) OnSubMessage(channel string, msg string) {
	if this.channel == channel {
		infos := strings.Split(msg, "#")
		if len(infos) == 2 {
			if infos[0] == GONODE_NEW_NODE {
				this.callbacks.OnNewNode(infos[1])
			}
		}
	}
}
func (this *RedisDC) OnSubError(channel string, err error) {
	if this.channel == channel {
		this.callbacks.OnDCError(err)
	}
}
