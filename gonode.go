package gonode

import (
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/itfantasy/gonode/behaviors/gen_monitor"
	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/behaviors/supervisor"

	"github.com/itfantasy/gonode/core/datacenter"
	"github.com/itfantasy/gonode/core/erl"
	"github.com/itfantasy/gonode/core/logger"
	"github.com/itfantasy/gonode/core/sysmonitor"

	"github.com/itfantasy/gonode/nets"
	"github.com/itfantasy/gonode/nets/kcp"
	"github.com/itfantasy/gonode/nets/tcp"
	"github.com/itfantasy/gonode/nets/ws"
)

type GoNode struct {
	info     *gen_server.NodeInfo
	behavior gen_server.GenServer

	logger *logger.Logger
	dc     datacenter.IDataCenter

	monitor    gen_monitor.GenMonitor
	monitoring *sysmonitor.SysMonitoring

	super supervisor.Supervisor

	netWorkers map[string]nets.INetWorker

	event *EventHandler
	lock  sync.RWMutex
}

// -------------- init ----------------

func (g *GoNode) Bind(behavior interface{}) error {
	switch behavior.(type) {
	case gen_server.GenServer:
		g.behavior = behavior.(gen_server.GenServer)
	case gen_monitor.GenMonitor:
		g.monitor = behavior.(gen_monitor.GenMonitor)
	case supervisor.Supervisor:
		super := behavior.(supervisor.Supervisor)
		superNode := supervisor.NewSuperNode()
		err := superNode.BindAndInit(SUPERVISOR, super, ALLNODES, CHAN_LOG, CHAN_MONI)
		if err != nil {
			return errors.New("Bind Supervisor Failed!!" + err.Error())
		}
		g.super = super
		g.behavior = superNode
	default:
		return errors.New("illegal behavior type!!")
	}
	return nil
}

func (g *GoNode) Launch() {
	defer g.onDispose()

	// mandatory multicore CPU enabled
	runtime.GOMAXPROCS(runtime.NumCPU())

	// get the node self info config
	if g.behavior == nil {
		fmt.Println("Initialize Faild!! You must bind a server behavior at first!!")
		return
	}
	info := g.behavior.Setup()
	if info == nil {
		fmt.Println("Initialize Faild!! Can not setup an correct nodeinfo!!")
		return
	}
	g.info = info
	g.event = newEventHandler(g)
	erl.BindErrorDigester(g.event)

	// init the logger
	logger, warn := logger.NewLogger(g.info.Id, g.info.LogLevel, CHAN_LOG, g.info.LogComp)
	if warn != nil {
		fmt.Println("Warning!! Can not create the Component for Logger, we will use the default Console Logger!" + warn.Error())
	}
	g.logger = logger

	// init the dc
	dc, err := datacenter.NewDataCenter(g.info.RegComp)
	if err != nil {
		fmt.Println("Initialize Faild!! Init the DataCenter failed!!" + err.Error())
		return
	}
	g.dc = dc
	g.dc.BindCallbacks(g.event)
	err2 := g.dc.RegisterAndDetect(g.info, CHAN_REG, 5000)
	if err2 != nil {
		fmt.Println("Initialize Faild!! Register to the DataCenter failed!!" + err2.Error())
		return
	}

	// init the monitor
	if g.monitor != nil {
		monitoring, err := sysmonitor.NewSysMonitoring(g.info.Id, g.monitor, CHAN_MONI)
		if err != nil {
			fmt.Println("Initialize Faild!! You have binded a event behavior, but the eventconf is incorrect!!" + err.Error())
			return
		}
		g.monitoring = monitoring
		g.monitoring.StartMonitoring()
	}

	theUrl, err := g.getListenUrl(g.info.Url)
	if err != nil {
		fmt.Println("Initialize Faild!! Can not parse the url!!")
		g.logger.Error(err.Error())
	}
	g.Listen(theUrl)

	fmt.Println(` ------- itfantasy.github.io -------
   ______      _   __          __   
  / ____/___  / | / /___  ____/ /__ 
 / / __/ __ \/  |/ / __ \/ __  / _ \
/ /_/ / /_/ / /|  / /_/ / /_/ /  __/
\____/\____/_/ |_/\____/\__,_/\___/ 

 --------------------------------- ` + VERSION)

	fmt.Println(g.info.ToString())
	g.logger.Info("node is starting... " + g.info.Id)
	g.behavior.Start()
	select {}
	g.logger.Error("shuting down!!!")
}

// -------------- props ------------------

func (g *GoNode) Info() *gen_server.NodeInfo {
	return g.info
}

func (g *GoNode) Self() string {
	return g.info.Id
}

func (g *GoNode) Origin() string {
	return nets.CombineOriginInfo(g.info.Id, g.info.Url, g.info.Sig)
}

func (g *GoNode) Logger() *logger.Logger {
	return g.logger
}

// -------------- net ------------------

func (g *GoNode) getListenUrl(url string) (string, error) {
	infos := strings.Split(url, "://") // get the header of protocol
	if len(infos) != 2 {
		return "", errors.New("illegal url!" + url)
	}
	proto := infos[0]
	ipAndPort := strings.Split(infos[1], ":")
	if len(ipAndPort) != 2 {
		return "", errors.New("illegal url!" + url)
	}
	if !g.info.Pub {
		return g.info.Url, nil
	}
	return proto + "://" + "0.0.0.0" + ":" + ipAndPort[1], nil
}

func (g *GoNode) netWorker(url string) nets.INetWorker {
	if g.netWorkers == nil {
		g.netWorkers = make(map[string]nets.INetWorker)
	}
	_, exists := g.netWorkers[url]
	if !exists {
		proto := strings.Split(url, "://")[0] // get the header of protocol
		switch proto {
		case (string)(nets.WS):
			g.netWorkers[url] = ws.NewWSNetWorker()
		case (string)(nets.KCP):
			g.netWorkers[url] = kcp.NewKcpNetWorker()
		case (string)(nets.TCP):
			g.netWorkers[url] = tcp.NewTcpNetWorker()
		}
		g.netWorkers[url].BindEventListener(g.event)
	}
	return g.netWorkers[url]
}

func (g *GoNode) Listen(url string) {
	go func() {
		err := g.netWorker(url).Listen(url)
		if err != nil {
			g.logger.Error(err.Error())
			g.onError(g.info.Id, err)
		}
	}()
}

func (g *GoNode) Connnect(nickid string, url string) error {
	exist := nets.IsIdExists(nickid)
	if exist {
		g.logger.Info("there is a same id in local record:" + url + "#" + nickid)
		return nil
	}
	return g.netWorker(url).Connect(nickid, url, g.Origin())
}

func (g *GoNode) Send(id string, msg []byte) error {
	conn, _, netWorker, exist := nets.GetInfoConnById(id)
	if !exist {
		return errors.New("there is not the id in local record!")
	}
	return netWorker.Send(conn, msg)
}

func (g *GoNode) Close(id string) error {
	conn, _, netWorker, exist := nets.GetInfoConnById(id)
	if !exist {
		return errors.New("there is not the id in local record!")
	}
	return netWorker.Close(id, conn)
}

func (g *GoNode) checkTargetId(id string) bool {
	if g.info.BackEnds == ALLNODES && id != g.info.Id {
		return true
	}

	label := nets.Label(id)
	backEnds := strings.Split(g.info.BackEnds, ",")
	for _, item := range backEnds {
		if item == label {
			return true
		}
	}
	return false
}

// -------------- other ----------------

func (g *GoNode) reportError(err interface{}) {
	title := "!!! Auto Recovering..."
	content := fmt.Sprint(err) +
		"\r=============== - CallStackInfo - =============== \r" + string(debug.Stack())
	if g.monitor != nil {
		g.monitor.OnReportError(g.info.Id, title, content)
	}
	g.logger.Error(title + content)
}
