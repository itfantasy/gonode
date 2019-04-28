package gonode

import (
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/core/datacenter"
	"github.com/itfantasy/gonode/core/logger"
	"github.com/itfantasy/gonode/nets"
	"github.com/itfantasy/gonode/nets/kcp"
	"github.com/itfantasy/gonode/nets/tcp"
	"github.com/itfantasy/gonode/nets/ws"
	"github.com/itfantasy/gonode/utils/snowflake"

	log "github.com/jeanphorn/log4go"
)

type GoNode struct {
	info     *gen_server.NodeInfo
	behavior gen_server.GenServer
	event    *EventHandler

	logger *log.Filter
	dc     datacenter.IDataCenter

	netWorkers map[string]nets.INetWorker

	lock sync.RWMutex
}

// -------------- init ----------------

func (this *GoNode) Initialize(behavior gen_server.GenServer) {

	// mandatory multicore CPU enabled
	runtime.GOMAXPROCS(runtime.NumCPU())
	nets.InitKvvk()

	// get the node self info config
	this.behavior = behavior
	info := this.behavior.Setup()
	if info == nil {
		fmt.Println("Initialize Faild!! Can not setup an correct nodeinfo!!")
		return
	}
	this.info = info
	this.event = newEventHandler(this)

	// init the logger
	logger, warn := logger.NewLogger(this.info.Id, this.info.LogLevel, GONODE_LOG_CHAN, this.info.LogComp)
	if warn != nil {
		this.logger.Warn(warn.Error())
		fmt.Println("Warning!! Can not create the Component for Logger, we will use the default Console Logger!")
	}
	this.logger = logger

	// init the dc
	dc, err := datacenter.NewDataCenter(this.info.RegComp)
	if err != nil {
		fmt.Println("Initialize Faild!! Init the DataCenter failed!!")
		this.logger.Error(err.Error())
		return
	}
	this.dc = dc
	this.dc.BindCallbacks(this.event)
	err2 := this.dc.RegisterAndDetect(this.info, GONODE_REG_CHAN, 5000)
	if err2 != nil {
		fmt.Println("Initialize Faild!! Register to the DataCenter failed!!")
		this.logger.Error(err2.Error())
	}

	theUrl, err := this.getListenUrl(this.info.Url)
	if err != nil {
		fmt.Println("Initialize Faild!! Can not parse the url!!")
		this.logger.Error(err.Error())
	}
	this.Listen(theUrl)
}

func (this *GoNode) Bind(behavior gen_server.GenServer) {
	this.behavior = behavior
}

func (this *GoNode) Sync() {
	defer this.onDispose()
	this.logger.Info("node starting... " + this.info.Id)
	this.behavior.Start()
	select {}
	this.logger.Error("shuting down!!!")
}

// -------------- props ------------------

func (this *GoNode) Info() *gen_server.NodeInfo {
	return this.info
}

func (this *GoNode) Self() string {
	return this.info.Id
}

func (this *GoNode) Origin() string {
	return nets.CombineOriginInfo(this.info.Id, this.info.Url, this.info.Sig)
}

func (this *GoNode) Logger() *log.Filter {
	return this.logger
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
		case (string)(nets.TCP):
			this.netWorkers[url] = new(tcp.TcpNetWorker)
			break
		}
		this.netWorkers[url].BindEventListener(this.event)
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
			this.onError(this.info.Id, err)
		}
	}()
}

func (this *GoNode) Connnect(nickid string, url string) error {
	return this.netWorker(url).Connect(nickid, url, this.Origin())
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

func (this *GoNode) randomCntId() string {
	return "cnt-" + snowflake.Generate()
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

// -------------- other ----------------

func (this *GoNode) autoRecover() {
	err := recover()
	if err != nil {
		this.logger.Error("auto recovering..." + fmt.Sprint(err))
		debug.PrintStack()
	}
}
