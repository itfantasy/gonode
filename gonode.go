package gonode

import (
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/itfantasy/gonode/behaviors/cmd"
	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/components"
	"github.com/itfantasy/gonode/components/logger"
	"github.com/itfantasy/gonode/components/redis"
	"github.com/itfantasy/gonode/nets"
	"github.com/itfantasy/gonode/nets/kcp"
	"github.com/itfantasy/gonode/nets/ws"
	"github.com/itfantasy/gonode/utils/crypt"
	"github.com/itfantasy/gonode/utils/json"
	"github.com/itfantasy/gonode/utils/timer"

	log "github.com/jeanphorn/log4go"
)

type GoNode struct {
	info     *gen_server.NodeInfo
	behavior gen_server.GenServer
	logger   *log.Filter
	logcomp  components.IComponent

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

func Logger() *log.Filter {
	return Node().Logger()
}

func Log(msg string) {
	Logger().Debug(msg)
}

func Console(obj interface{}) {
	msg, err := json.Encode(obj)
	if err != nil {
		Log(msg)
	}
}

func Error(msg string) {
	Logger().Error(msg)
}

// -------------- init ----------------

func (this *GoNode) Bind(behavior gen_server.GenServer) {
	this.behavior = behavior
}

func (this *GoNode) Initialize(behavior gen_server.GenServer) {

	defer this.Dispose()

	// mandatory multicore CPU enabled
	runtime.GOMAXPROCS(runtime.NumCPU())
	// get the node self info config
	this.behavior = behavior
	info := this.behavior.Setup()
	if info == nil {
		fmt.Println("Initialize Faild!! Can not setup an correct nodeinfo!!")
		return
	}
	this.info = info
	// init the logger
	if this.info.LogComp != "" {
		logcomp, err := components.NewComponent(this.info.LogComp)
		if err != nil {
			fmt.Println("Warning!! Can not create the Component for Logger, we will use the default Console Logger!")
		}
		this.logcomp = logcomp
	}

	this.logger = logger.NewLogger(this.info.Id, this.info.LogLevel, GONODE_LOG_CHAN, this.logcomp)

	// init the core redis
	regcomp, err := components.NewComponent(this.info.RegComp)
	if err != nil {
		fmt.Println("Initialize Faild!! Can not create the Core Register Component!!")
		this.logger.Error(err.Error())
		return
	}
	this.coreRedis = regcomp.(*redis.Redis)
	// sub the redis channel
	this.coreRedis.BindSubscriber(this)
	go this.coreRedis.Subscribe(GONODE_PUB_CHAN)

	theUrl, err := this.getListenUrl(this.info.Url)
	if err != nil {
		fmt.Println("Initialize Faild!! Can not parse the url!!")
		this.logger.Error(err.Error())
	}

	// register self info to core redis
	this.registerSelf()
	this.Listen(theUrl)

	// check if auto detect
	if this.info.BackEnds != "" {
		go this.autoDetect()
	}

	this.behavior.Start()
	select {}
	this.logger.Error("shuting down!!!")
}

func (this *GoNode) Dispose() {

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
	this.logger.Error(err.Error())
}

// -------------- net ------------------

func (this *GoNode) getListenUrl(url string) (string, error) {
	infos := strings.Split(url, "://") // get the header of protocol
	if len(infos) != 2 {
		return "", errors.New("illegal url!" + url)
	}
	proto := infos[0]
	ipAndPort := strings.Split(infos[1], ":")
	if len(ipAndPort) != 2 {
		return "", errors.New("illegal url!" + url)
	}
	if !this.info.Pub {
		return this.info.Url, nil
	}
	return proto + "://" + "0.0.0.0" + ":" + ipAndPort[1], nil
}

func (this *GoNode) netWorker(url string) nets.INetWorker {
	if this.netWorkers == nil {
		this.netWorkers = make(map[string]nets.INetWorker)
	}
	infos := strings.Split(url, "://") // get the header of protocol
	proto := infos[0]
	_, exists := this.netWorkers[url]
	if !exists {
		switch proto {
		case (string)(nets.WS):
			this.netWorkers[url] = new(ws.WSNetWorker)
			break
		case (string)(nets.KCP):
			this.netWorkers[url] = new(kcp.KcpNetWorker)
			break
		}
		this.netWorkers[url].BindEventListener(this)
	} else {
		this.logger.Warn("the url has been listening!" + url)
	}
	return this.netWorkers[url]
}

func (this *GoNode) Listen(url string) {
	go func() {
		err := this.netWorker(url).Listen(url)
		if err != nil {
			this.logger.Error(err.Error())
			this.OnError(this.info.Id, err)
		}
	}()
}

func (this *GoNode) Connnect(url string, origin string) error {
	return this.netWorker(url).Connect(url, origin)
}

func (this *GoNode) Send(id string, msg []byte) error {
	conn, _, netWorker, exist := nets.GetInfoConnById(id)
	if !exist {
		return errors.New("there is not the id in local record!")
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
		if !this.info.Pub {
			this.logger.Info("not a inside node! give up the url:" + url)
			return "", false
		} else {
			connId := this.randomCntId()
			return connId, true
		}
	} else {
		exist := nets.IsIdExists(info.Id)
		if exist {
			this.logger.Info("there is a same id in local record:" + url + "|" + info.Id)
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
		if this.checkTargetId(id) {
			this.logger.Info("a new node has been found!", id)
			// find the node url by the id
			url, err := this.getNodeUrlById(id)
			if err == nil {
				err2 := this.Connnect(url, this.info.Url)
				if err2 != nil {
					this.logger.Error(err2.Error() + "[" + id + "]")
				}
			} else {
				this.logger.Error(err.Error() + "[" + id + "]")
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
		this.logger.Error(err.Error())

	}

	this.coreRedis.Set("gonode_"+this.info.Id, this.info.Url)
	this.coreRedis.Set(this.info.Url, string(infoStr))
	this.coreRedis.SAdd(GONODE_INFO, this.info.Id)

	msg := cmd.NewNode(this.info.Id)
	this.PublishMsg(msg)

	this.logger.Info("report the node info:" + msg)
}

func (this *GoNode) randomCntId() string {
	return "cnt@" + this.info.Id + "$" + crypt.C32()
}

func (this *GoNode) checkTargetId(id string) bool {
	label := strings.Split(id, "-")[0]
	backEnds := strings.Split(this.info.BackEnds, ",")
	for _, item := range backEnds {
		if item == label {
			return true
		}
	}
	return false
}

// -------------- redis --------------------

func (this *GoNode) CoreRedis() *redis.Redis {
	return this.coreRedis
}

// -------------- logger -------------------

func (this *GoNode) Logger() *log.Filter {
	return this.logger
}

// -------------- other ----------------

func (this *GoNode) autoRecover() {
	err := recover()
	if err != nil {
		this.logger.Error("auto recovering..." + fmt.Sprint(err))
		debug.PrintStack()
	}
}
