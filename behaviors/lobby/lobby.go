package lobby

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/utils/ini"
	"github.com/itfantasy/gonode/utils/io"
	//	"github.com/itfantasy/gonode/utils/timer"
)

type Lobby struct {
	server LobbyServer
}

func (this *Lobby) SelfNodeInfo() (*gen_server.NodeInfo, error) {
	conf, err := ini.Load(io.CurDir() + "conf.ini")
	if err != nil {
		return nil, err
	}
	nodeInfo := new(gen_server.NodeInfo)
	nodeInfo.Tag = conf.Get("node", "tag")
	nodeInfo.Id = conf.Get("node", "id")
	nodeInfo.Url = conf.Get("node", "url")
	nodeInfo.RedUrl = conf.Get("redis", "url")
	nodeInfo.RedPool = conf.GetInt("redis", "pool", 0)
	nodeInfo.RedDB = conf.GetInt("redis", "db", 0)
	nodeInfo.RedAuth = conf.Get("redis", "auth")
	nodeInfo.AutoDetect = conf.GetInt("net", "autodetect", 0) > 0
	nodeInfo.Net = conf.Get("net", "net")
	return nodeInfo, nil
}
func (this *Lobby) IsInterestedIn(string) bool {
	return false
}
func (this *Lobby) Start() {
	fmt.Println("node starting...")
	this.server.Start()
}
func (this *Lobby) Update() {

}
func (this *Lobby) OnConn(id string) {

}
func (this *Lobby) OnMsg(id string, msg []byte) {
	if strings.Contains(id, "room") {
		// native logic for roomserver
	} else {
		this.server.OnMsg(id, msg)
	}
}
func (this *Lobby) OnClose(id string) {

}
func (this *Lobby) OnShell(id string, msg string) {

}
func (this *Lobby) OnReload(tag string) error {
	return nil
}
func (this *Lobby) CreateConnId() string {
	return "cnt" + strconv.Itoa(rand.Intn(100000))
}
func (this *Lobby) Initialize(server LobbyServer) {
	this.server = server
	gonode.Node().Initialize(this)
}
