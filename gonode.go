package gonode

import (
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/components"
	"github.com/itfantasy/gonode/components/email"
	"github.com/itfantasy/gonode/core/datacenter"
	"github.com/itfantasy/gonode/core/erl"
	"github.com/itfantasy/gonode/core/logger"
	"github.com/itfantasy/gonode/nets"
	"github.com/itfantasy/gonode/nets/kcp"
	"github.com/itfantasy/gonode/nets/tcp"
	"github.com/itfantasy/gonode/nets/ws"
	"github.com/itfantasy/gonode/utils/snowflake"
)

type GoNode struct {
	info     *gen_server.NodeInfo
	behavior gen_server.GenServer
	event    *EventHandler

	logger *logger.Logger
	dc     datacenter.IDataCenter

	mail *email.Email

	netWorkers map[string]nets.INetWorker

	lock sync.RWMutex
}

// -------------- init ----------------

func (g *GoNode) Bind(behavior gen_server.GenServer) {
	g.behavior = behavior
}

func (g *GoNode) Launch() {
	defer g.onDispose()

	// mandatory multicore CPU enabled
	runtime.GOMAXPROCS(runtime.NumCPU())
	nets.InitKvvk()

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
	erl.BindErrorReporter(g.event)

	// init the logger
	logger, warn := logger.NewLogger(g.info.Id, g.info.LogLevel, CHAN_LOG, g.info.LogComp)
	if warn != nil {
		g.logger.Warn(warn.Error())
		fmt.Println("Warning!! Can not create the Component for Logger, we will use the default Console Logger!")
	}
	g.logger = logger

	// init the email
	if g.info.RepComp != "" {
		repComp, err := components.NewComponent(g.info.RepComp)
		if err != nil {
			g.logger.Error(err.Error())
			fmt.Println("Initialize Faild!! Have setted up the Error Reporter, but canot create it as an Email Componment! .." + err.Error())
			return
		}
		mail, ok := repComp.(*email.Email)
		if !ok {
			fmt.Println("Initialize Faild!! Have setted up the Error Reporter, but canot create it as an Email Componment! ..")
		}
		g.mail = mail
	}

	// init the dc
	dc, err := datacenter.NewDataCenter(g.info.RegComp)
	if err != nil {
		fmt.Println("Initialize Faild!! Init the DataCenter failed!!")
		g.logger.Error(err.Error())
		return
	}
	g.dc = dc
	g.dc.BindCallbacks(g.event)
	err2 := g.dc.RegisterAndDetect(g.info, CHAN_REG, 5000)
	if err2 != nil {
		fmt.Println("Initialize Faild!! Register to the DataCenter failed!!")
		g.logger.Error(err2.Error())
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
	infos := strings.Split(url, "://") // get the header of protocol
	proto := infos[0]
	_, exists := g.netWorkers[url]
	if !exists {
		switch proto {
		case (string)(nets.WS):
			g.netWorkers[url] = new(ws.WSNetWorker)
			break
		case (string)(nets.KCP):
			g.netWorkers[url] = new(kcp.KcpNetWorker)
			break
		case (string)(nets.TCP):
			g.netWorkers[url] = new(tcp.TcpNetWorker)
			break
		}
		g.netWorkers[url].BindEventListener(g.event)
	} else {
		g.logger.Warn("the url has been listening!" + url)
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

func (g *GoNode) GetAllConnIds() []string {
	return nets.GetAllConnIds()
}

func (g *GoNode) randomCntId() string {
	return "cnt-" + snowflake.Generate()
}

func (g *GoNode) checkTargetId(id string) bool {
	label := strings.Split(id, "-")[0]
	backEnds := strings.Split(g.info.BackEnds, ",")
	for _, item := range backEnds {
		if item == label {
			return true
		}
	}
	return false
}

// -------------- other ----------------

func (g *GoNode) autoRecover() {
	if err := recover(); err != nil {
		g.reportError(err)
	}
}

func (g *GoNode) reportError(err interface{}) {
	title := "!!! Auto Recovering..."
	content := fmt.Sprint(err) +
		"\r=============== - CallStackInfo - =============== \r" + string(debug.Stack())
	if g.mail != nil {
		go func() {
			err := g.mail.SendTo(g.info.RepTo, title, strings.Replace(content, "\r", "<br/>", -1))
			if err != nil {
				g.logger.Error(err)
			}
		}()
	}
	g.logger.Error(title + content)
}
