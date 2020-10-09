package gonode

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/behaviors/logger"
	"github.com/itfantasy/gonode/behaviors/monitor"
	"github.com/itfantasy/gonode/behaviors/supervisor"
	"github.com/itfantasy/gonode/core/datacenter"
	"github.com/itfantasy/gonode/nets"
)

type GoNode struct {
	info       *gen_server.NodeInfo
	behavior   gen_server.GenServer
	logLevel   int
	logWriter  logger.LogWriter
	logger     *logger.Logger
	monitor    monitor.GenMonitor
	monitoring *monitor.SysMonitoring
	super      supervisor.Supervisor
	dc         datacenter.IDataCenter
	netWorkers map[string]nets.INetWorker
	event      *EventHandler
	lock       sync.RWMutex
}

// -------------- init ----------------

func (g *GoNode) Bind(behavior gen_server.GenServer) {
	g.behavior = behavior
}

func (g *GoNode) BindMonitor(monitor monitor.GenMonitor) {
	g.monitor = monitor
}

func (g *GoNode) BindLogger(logWriter logger.LogWriter, logLevel int) {
	g.logWriter = logWriter
	g.logLevel = logLevel
}

func (g *GoNode) Supervise(super supervisor.Supervisor) error {
	if g.behavior != nil {
		return errors.New("Bind Supervisor Failed!! This node has been binded by gen_server or supervisor")
	}
	superNode := supervisor.NewSuperNode()
	err := superNode.InitSupervisor(super)
	if err != nil {
		return errors.New("Bind Supervisor Failed!!" + err.Error())
	}
	g.super = super
	g.behavior = superNode
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

	// init the logger
	logger := logger.NewLogger(g.info.NodeId, g.logLevel, g.logWriter)
	g.logger = logger

	// init the dc
	dc, err := datacenter.NewDataCenter(g.info.RegDC)
	if err != nil {
		fmt.Println("Initialize Faild!! Init the DataCenter failed!!" + err.Error())
		return
	}
	g.dc = dc
	g.dc.BindCallbacks(g.event)
	err2 := g.dc.RegisterAndDetect(g.info, g.info.NameSpace, 5000)
	if err2 != nil {
		fmt.Println("Initialize Faild!! Register to the DataCenter failed!!" + err2.Error())
		return
	}

	// init the monitor
	if g.monitor != nil {
		monitoring, err := monitor.NewSysMonitoring(g.info.NodeId, g.monitor)
		if err != nil {
			fmt.Println("Initialize Faild!! You have binded a event behavior, but the eventconf is incorrect!!" + err.Error())
			return
		}
		g.monitoring = monitoring
		g.monitoring.StartMonitoring()
	}

	if len(g.info.EndPoints) > 0 {
		for _, endPoint := range g.info.EndPoints {
			theUrl, err := g.getListenUrl(endPoint)
			if err != nil {
				fmt.Println("Initialize Faild!! Can not parse the url!!")
				g.logger.Error(err.Error())
			} else {
				g.Listen(theUrl)
			}
		}
	}

	fmt.Println(` ------- itfantasy.github.io -------
   ______      _   __          __   
  / ____/___  / | / /___  ____/ /__ 
 / / __/ __ \/  |/ / __ \/ __  / _ \
/ /_/ / /_/ / /|  / /_/ / /_/ /  __/
\____/\____/_/ |_/\____/\__,_/\___/ 

 -----------:: gonode ::---------- ` + VERSION)

	fmt.Println(g.info.ToString())
	g.behavior.Start()
	g.logger.Info("The node has been Launched!" + g.info.NodeId)
	select {}
	g.logger.Error("shuting down!!!")
}

// -------------- props ------------------

func (g *GoNode) Info() *gen_server.NodeInfo {
	return g.info
}

func (g *GoNode) Self() string {
	return g.info.NodeId
}

func (g *GoNode) Origin() string {
	return nets.CombineOriginInfo(g.info.NodeId, g.info.EndPoints[0], g.info.Sig)
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
	gridEvn := os.Getenv("GRID_NODE_ID")
	if gridEvn != "" && g.info.IsPub {
		return proto + "://" + "0.0.0.0" + ":" + ipAndPort[1], nil
	}
	return url, nil
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
			g.netWorkers[url] = nets.NewWSNetWorker()
		case (string)(nets.KCP):
			g.netWorkers[url] = nets.NewKcpNetWorker()
		case (string)(nets.TCP):
			g.netWorkers[url] = nets.NewTcpNetWorker()
		}
		g.netWorkers[url].BindEventListener(g.event)
	}
	return g.netWorkers[url]
}

func (g *GoNode) Listen(url string) {
	go func() {
		g.lock.Lock()
		netWorker := g.netWorker(url)
		g.lock.Unlock()
		if err := netWorker.Listen(url); err != nil {
			g.logger.Error(err.Error())
			g.onError(g.info.NodeId, err)
		}
	}()
}

func (g *GoNode) Connnect(nodeId string, url string) error {
	exist := nets.NodeConned(nodeId)
	if exist {
		g.logger.Info("The nickid has been existed!" + url + "#" + nodeId)
		return nil
	}
	return g.netWorker(url).Connect(nodeId, url, g.Origin())
}

func (g *GoNode) Send(id string, msg []byte) error {
	conn, _, netWorker, exist := nets.GetInfoConnById(id)
	if !exist {
		return errors.New("Cannot find the sending id!" + id)
	}
	return netWorker.Send(conn, msg)
}

func (g *GoNode) SendAll(ids []string, msg []byte) []error {
	errs := make([]error, 0, len(ids))
	for _, id := range ids {
		err := g.Send(id, msg)
		errs = append(errs, err)
	}
	return errs
}

func (g *GoNode) Close(id string) error {
	conn, _, netWorker, exist := nets.GetInfoConnById(id)
	if !exist {
		return errors.New("Cannot find the closing id!" + id)
	}
	return netWorker.Close(id, conn)
}

func (g *GoNode) checkTargetId(id string) bool {
	if g.info.BackEnds == ALLNODES && id != g.info.NodeId {
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
		g.monitor.OnReportError(g.info.NodeId, title, content)
	}
	g.logger.Error(title + content)
}
