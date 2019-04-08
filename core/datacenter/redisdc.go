package datacenter

import (
	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/components/redis"
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
func (this *RedisDC) RegisterAndDetect(info *gen_server.NodeInfo, channel string, msfrequency int) {
	this.info = info
	this.channel = channel

}
func (this *RedisDC) GetNodeInfo(id string) (*gen_server.NodeInfo, error) {
	if this.info.Id == id {
		return this.info, nil
	}

	return nil, nil
}
func (this *RedisDC) CheckNode(id string, origin string) bool {
	return false
}

func (this *RedisDC) OnSubscribe(channel string) {

}
func (this *RedisDC) OnSubMessage(channel string, msg string) {

}
func (this *RedisDC) OnSubError(channel string, err error) {

}
