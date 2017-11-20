package gonode

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"

	native_redis "github.com/garyburd/redigo/redis"

	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/components/logger"
	"github.com/itfantasy/gonode/components/redis"
	"github.com/itfantasy/gonode/nets"
	"github.com/itfantasy/gonode/nets/ws"
	"github.com/itfantasy/gonode/roles/shellcmd/cmd"
	"github.com/itfantasy/gonode/utils/json"
	"github.com/itfantasy/gonode/utils/timer"
)

type GoNode struct {
	info      *gen_server.NodeInfo
	behavior  gen_server.GenServer
	logger    *logger.Logger
	coreRedis *redis.Redis
	netWorker nets.INetWorker
	lock      sync.RWMutex
}

var node *GoNode = nil

func Node() *GoNode {
	if node == nil {
		node = &GoNode{}
	}
	return node
}

// -------------- init ----------------

func (this *GoNode) Initialize(behavior gen_server.GenServer) {
	// mandatory multicore CPU enabled
	runtime.GOMAXPROCS(runtime.NumCPU())
	// init the logger
	this.logger = new(logger.Logger)
	// get the node self info config
	this.info = this.behavior.SelfNodeInfo()
	// init the core redis
	this.coreRedis = new(redis.Redis)
	if err := this.coreRedis.Conn(this.info.RedCore,
		this.info.RedDB,
		this.info.RedPool,
		this.info.RedAuth); err != nil {
		this.logger.Error(this.sprinfLog("cannot connect to the core redis!!"))
	}
	// sub the redis channel
	this.coreRedis.Subscribe(GONODE_PUB_CHAN)
	go this.handleSubscribe()

	if this.info.Net != "" {
		// init the networker
		this.initNetWorker()
		// register self info to core redis
		this.registerSelf()
		go this.netWorker.Listen(this.info.Url)
	}

	// check if auto detect
	if this.info.AutoDetect {
		this.autoDetect()
	}

	this.behavior.Start()

	for {
		timer.Sleep(1)
		this.behavior.Update()
	}
}

func (this *GoNode) initNetWorker() {
	url := this.info.Url
	infos := strings.Split(url, "://") // get the header of protocol
	switch infos[0] {
	case (string)(nets.WS):
		this.netWorker = new(ws.WSNetWorker)
	}
	this.netWorker.BindEventListener(this)
}

// -------------- redis pub/sub ------------------

const (
	// pub channel
	GONODE_PUB_CHAN string = "GONODE_PUB_CHAN"
	// log channel
	GONODE_LOG_CHAN string = "GONODE_LOG_CHAN"
	// all nodes infos
	GONODE_INFO string = "GONODE_INFO"
)

func (this *GoNode) PublishMsg(msg string) {
	this.coreRedis.Publish(GONODE_PUB_CHAN, msg)
}

func (this *GoNode) handleSubscribe() {
	for {
		switch v := this.coreRedis.Psc.Receive().(type) {
		case native_redis.Message:
			this.onShell(v.Channel, string(v.Data))
		case native_redis.Subscription:
			this.logger.Info(fmt.Sprintf("%s: %s %d\n", v.Channel, v.Kind, v.Count))
		case error:
			this.logger.Error(this.sprinfLog(v.Error()))
			return
		}
	}
}

// -------------- net ------------------

func (node *GoNode) getNodeInfo(url string) (*gen_server.NodeInfo, error) {
	infoStr, err := node.coreRedis.Get(url)
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

func (this *GoNode) CheckUrlLegal(url string) (string, bool) {
	// find the node info by redis at first
	info, err := node.getNodeInfo(url)
	if err != nil {
		// cannot find the node in lan
		if node.info.Net == "LAN" {
			this.logger.Info(this.sprinfLog("can not find the node! give up the url:" + url))
			return "", false
		} else { // if node.Info.Net == "WAN"
			connId := this.behavior.CreateConnId()
			return connId, true
		}
	} else {
		exist := this.netWorker.IsIdExists(info.Id)
		if exist {
			this.logger.Info(this.sprinfLog("there is a same id in local record:" + url + "/" + info.Id))
			return "", false
		}
		return info.Id, true
	}
}

func (this *GoNode) autoDetect() {
	this.logger.Info("auto detecting other nodes..")
	for {
		timer.Sleep(5000)
		// get all node infos from the coreRedis and compare with the local record
		ids, err := this.coreRedis.SMembers(GONODE_INFO)
		if err != nil {
			this.logger.Error("cannot find the nodes list from redis!!" + err.Error())
			continue
		}
		for _, id := range ids {
			this.checkNewNode(id)
		}
	}
}

// when a new node is found
func (this *GoNode) checkNewNode(id string) {
	this.lock.Lock()
	defer this.lock.Unlock()

	exist := this.netWorker.IsIdExists(id)
	if !exist {
		// check the local node is interested in the new node
		if this.behavior.IsInterestedIn(id) {
			this.logger.Info(this.sprinfLog("a new node has been found![" + id + "]"))
			// find the node url by the id
			url, err := this.getNodeUrlById(id)
			if err == nil {
				err2 := this.netWorker.Connect(url, this.info.Url)
				if err2 != nil {
					this.logger.Error(this.sprinfLog(err2.Error() + "[" + id + "]"))
				}
			} else {
				this.logger.Error(this.sprinfLog(err.Error() + "[" + id + "]"))
			}
		}
	}
}

func (this *GoNode) getNodeUrlById(id string) (string, error) {
	url, err := node.coreRedis.Get("gonode_" + id)
	return url, err
}

func (this *GoNode) registerSelf() {
	info := this.behavior.SelfNodeInfo()
	infoStr, err := json.Encode(info)
	if err != nil {
		this.logger.Error(this.sprinfLog(err.Error()))
	}

	this.coreRedis.Set("gonode_"+info.Id, info.Url)
	this.coreRedis.Set(info.Url, string(infoStr))
	this.coreRedis.SAdd(GONODE_INFO, info.Id)

	msg := cmd.NewNode(info.Id)
	this.PublishMsg(msg)

	this.logger.Info(this.sprinfLog("report the node info:" + msg))
}

// -------------- redis --------------------

func (this *GoNode) CoreRedis() *redis.Redis {
	return this.coreRedis
}

// -------------- logger -------------------

func (this *GoNode) Logger() *logger.Logger {
	return this.logger
}

func (this *GoNode) ReportLog(msg string) {
	// report the log by the redis comp
	this.coreRedis.Publish(GONODE_LOG_CHAN, msg)
}

func (this *GoNode) sprinfLog(log string) string {
	return fmt.Sprintf("[%s]--%s", this.info.Id, log)
}

// -------------- other ----------------

func (this *GoNode) autoRecover() {
	err := recover()
	if err != nil {
		this.logger.Error("auto recovering..." + fmt.Sprint(err))
		debug.PrintStack()
	}
}
