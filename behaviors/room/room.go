package room

import (
	"fmt"
	"strings"

	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/utils/ini"
	"github.com/itfantasy/gonode/utils/io"
	//	"github.com/itfantasy/gonode/utils/timer"
)

type Room struct {
}

func (this *Room) SelfNodeInfo() (*gen_server.NodeInfo, error) {
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
func (this *Room) IsInterestedIn(id string) bool {
	if id == "lobby" { // the room will auto find the lobby, and try to build a conn to the lobby
		return true
	}
	return false
}
func (this *Room) Start() {
	fmt.Println("node starting...")
}
func (this *Room) Update() {

}
func (this *Room) OnConn(id string) {

}
func (this *Room) OnMsg(id string, msg []byte) {

}
func (this *Room) OnClose(id string) {

}
func (this *Room) OnShell(id string, msg string) {

}
func (this *Room) OnReload(tag string) error {
	return nil
}
func (this *Room) CreateConnId() string {
	return ""
}
func main() {
	room := new(Room)
	gonode.Node().Initialize(room)
}
