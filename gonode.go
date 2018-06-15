package gonode

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"

	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/behaviors/shellcmd/cmd"
	"github.com/itfantasy/gonode/components/logger"
	"github.com/itfantasy/gonode/components/redis"
	"github.com/itfantasy/gonode/nets"
	"github.com/itfantasy/gonode/nets/kcp"
	"github.com/itfantasy/gonode/nets/ws"
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

func Send(id string, msg []byte) {
	Node().NetWorker().Send(id, msg)
}

func Log(msg string) {
	Node().Logger().Info(Node().sprinfLog(msg))
}

func Console(obj interface{}) {
	msg, err := json.Encode(obj)
	if err != nil {
		Node().Logger().Debug(Node().sprinfLog(msg))
	}
}

func Error(msg string) {
	Node().Logger().Error(Node().sprinfLog(msg))
}

// -------------- init ----------------

func (this *GoNode) Initialize(behavior gen_server.GenServer) {
	// mandatory multicore CPU enabled
	runtime.GOMAXPROCS(runtime.NumCPU())
	// init the logger
	this.logger = new(logger.Logger)
	// get the node self info config
	this.behavior = behavior
	info, err := this.behavior.SelfNodeInfo()
	if err != nil {
		this.logger.Error("get the node self info err!" + err.Error())
		return
	}
	this.info = info
	// init the core redis
	this.coreRedis = redis.NewRedis()
	this.coreRedis.SetAuthor("", this.info.RedAuth)
	this.coreRedis.SetOption(redis.OPT_MAXPOOL, this.info.RedPool)
	if err := this.coreRedis.Conn(this.info.RedUrl, strconv.Itoa(this.info.RedDB)); err != nil {
		this.logger.Error(this.sprinfLog("cannot connect to the core redis!!"))
		return
	}
	// sub the redis channel
	this.coreRedis.BindSubscriber(this)
	go this.coreRedis.Subscribe(GONODE_PUB_CHAN)
	//this.handleSubscribe()

	if this.info.Net != "" {
		// init the networker
		if err := this.initNetWorker(); err != nil {
			this.logger.Error(this.sprinfLog(err.Error()))
			return
		}
		// register self info to core redis
		this.registerSelf()
		go this.netWorker.Listen(this.info.Url)
	}

	// check if auto detect
	if this.info.AutoDetect {
		go this.autoDetect()
	}

	this.behavior.Start()

	for {
		timer.Sleep(160)
		this.behavior.Update()
	}

	this.logger.Error("shuting down!!!")
}

func (this *GoNode) initNetWorker() error {
	url := this.info.Url
	infos := strings.Split(url, "://") // get the header of protocol
	switch infos[0] {
	case (string)(nets.WS):
		this.netWorker = new(ws.WSNetWorker)
		break
	case (string)(nets.KCP):
		this.netWorker = new(kcp.KcpNetWorker)
		break
	}
	return this.netWorker.BindEventListener(this)
}

// -------------- info --------------------
func (this *GoNode) Info() *gen_server.NodeInfo {
	return this.info
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

func (this *GoNode) OnSubscribe(channel string) {
	this.logger.Info("gonode has subscribed the channel:" + channel)
}

func (this *GoNode) OnSubMessage(channel string, msg string) {
	this.onShell(channel, msg)
}

func (this *GoNode) OnSubError(channel string, err error) {
	this.logger.Error(this.sprinfLog(err.Error()))
}

/*
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
*/

// -------------- net ------------------

func (this *GoNode) NetWorker() nets.INetWorker {
	return this.netWorker
}

func (this *GoNode) getNodeInfo(url string) (*gen_server.NodeInfo, error) {
	infoStr, err := this.coreRedis.Get(url)
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
			this.logger.Info(this.sprinfLog("there is a same id in local record:" + url + "|" + info.Id))
			return "", false
		}
		return info.Id, true
	}
}

func (this *GoNode) autoDetect() {
	for {
		timer.Sleep(5000)
		//this.logger.Info("auto detecting other nodes..")
		// get all node infos from the coreRedis and compare with the local record
		ids, err := this.coreRedis.SMembers(GONODE_INFO)
		if err != nil {
			this.logger.Error("cannot find the nodes list from redis!!" + err.Error())
			continue
		}
		for _, id := range ids {
			if id != this.info.Id {
				this.checkNewNode(id)
			}
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
	infoStr, err := json.Encode(this.info)
	if err != nil {
		this.logger.Error(this.sprinfLog(err.Error()))
	}

	this.coreRedis.Set("gonode_"+this.info.Id, this.info.Url)
	this.coreRedis.Set(this.info.Url, string(infoStr))
	this.coreRedis.SAdd(GONODE_INFO, this.info.Id)

	msg := cmd.NewNode(this.info.Id)
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
