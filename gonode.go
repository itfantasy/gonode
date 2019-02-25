package gonode

import (
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"

	"github.com/itfantasy/gonode/behaviors/cmd"
	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/components/logger"
	"github.com/itfantasy/gonode/components/redis"
	"github.com/itfantasy/gonode/nets"
	"github.com/itfantasy/gonode/nets/kcp"
	"github.com/itfantasy/gonode/nets/ws"
	"github.com/itfantasy/gonode/utils/json"
	"github.com/itfantasy/gonode/utils/timer"
)

type GoNode struct {
	info       *gen_server.NodeInfo
	behavior   gen_server.GenServer
	logger     *logger.Logger
	coreRedis  *redis.Redis
	netWorkers map[string]nets.INetWorker
	lock       sync.RWMutex
}

var node *GoNode = nil

func Node() *GoNode {
	if node == nil {
		node = &GoNode{}
	}
	return node
}

func Send(id string, msg []byte) {
	Node().Send(id, msg)
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

func (this *GoNode) Bind(behavior gen_server.GenServer) {
	this.behavior = behavior
}

func (this *GoNode) Initialize(behavior gen_server.GenServer) {
	// mandatory multicore CPU enabled
	runtime.GOMAXPROCS(runtime.NumCPU())
	// init the logger
	this.logger = new(logger.Logger)
	// get the node self info config
	this.behavior = behavior
	info, err := this.behavior.Setup()
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

	// register self info to core redis
	this.registerSelf()
	this.Listen(this.info.Url)

	// check if auto detect
	if this.info.AutoDetect {
		go this.autoDetect()
	}

	this.behavior.Start()
	select {}
	this.logger.Error("shuting down!!!")
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

// -------------- net ------------------

func (this *GoNode) netWorker(url string) nets.INetWorker {
	if this.netWorkers == nil {
		this.netWorkers = make(map[string]nets.INetWorker)
	}
	infos := strings.Split(url, "://") // get the header of protocol
	proto := infos[0]
	_, exists := this.netWorkers[proto]
	if !exists {
		switch proto {
		case (string)(nets.WS):
			this.netWorkers[proto] = new(ws.WSNetWorker)
			break
		case (string)(nets.KCP):
			this.netWorkers[proto] = new(kcp.KcpNetWorker)
			break
		}
		this.netWorkers[proto].BindEventListener(this)
	} else {
		this.logger.Warn(this.sprinfLog("there has been a same proto networker!" + proto))
	}
	return this.netWorkers[proto]
}

func (this *GoNode) Listen(url string) {
	go func() {
		err := this.netWorker(url).Listen(url)
		if err != nil {
			this.logger.Error(this.sprinfLog(err.Error()))
		}
	}()
}

func (this *GoNode) Connnect(url string, origin string) error {
	return this.netWorker(url).Connect(url, origin)
}

func (this *GoNode) Send(id string, msg []byte) error {
	conn, proto, exist := nets.GetInfoConnById(id)
	if !exist {
		return errors.New("there is not the id in local record!")
	}
	netWorker, exist := this.netWorkers[proto]
	if !exist {
		return errors.New("illegal proto!")
	}
	return netWorker.Send(conn, msg)
}

func (this *GoNode) GetAllConnIds() []string {
	return nets.GetAllConnIds()
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
	info, err := this.getNodeInfo(url)
	if err != nil {
		// cannot find the node in lan
		if !this.info.Public {
			this.logger.Info(this.sprinfLog("not a inside node! give up the url:" + url))
			return "", false
		} else {
			connId := this.behavior.OnRanId()
			return connId, true
		}
	} else {
		exist := nets.IsIdExists(info.Id)
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

	exist := nets.IsIdExists(id)
	if !exist {
		// check the local node is interested in the new node
		if this.behavior.OnDetect(id) {
			this.logger.Info(this.sprinfLog("a new node has been found![" + id + "]"))
			// find the node url by the id
			url, err := this.getNodeUrlById(id)
			if err == nil {
				err2 := this.Connnect(url, this.info.Url)
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
