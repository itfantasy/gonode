package lobby

import (
	"fmt"
	"strings"

	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/utils/ini"
	"github.com/itfantasy/gonode/utils/io"
	//	"github.com/itfantasy/gonode/utils/timer"
)

type Lobby struct {
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
}
func (this *Lobby) Update() {

}
func (this *Lobby) OnConn(id string) {

}
func (this *Lobby) OnMsg(id string, msg []byte) {
	if strings.Contains("room") {
		this.OnMsg4Room(id, msg)
	} else {
		this.OnMsg4Cnt(id, msg)
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
	return ""
}
func main() {
	lobby := new(Lobby)
	gonode.Node().Initialize(lobby)
}
