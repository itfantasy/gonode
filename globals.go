package gonode

import (
	"fmt"

	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/nets"
	"github.com/itfantasy/gonode/utils/json"
)

// -------------- global ----------------

var node *GoNode = nil

func init() {
	if node == nil {
		node = &GoNode{}
	}
}

func Launch() {
	node.Launch()
}

func Bind(behavior interface{}) error {
	return node.Bind(behavior)
}

func Listen(url string) {
	node.Listen(url)
}

func Connect(nickid string, url string) error {
	return node.Connnect(nickid, url)
}

func Send(id string, msg []byte) error {
	return node.Send(id, msg)
}

func Close(id string) error {
	return node.Close(id)
}

func Label(id string) string {
	return nets.Label(id)
}

func AllNodes() []string {
	return nets.AllNodes()
}

func Nodes(label string) []string {
	return nets.Nodes(label)
}

func AllPeers() []string {
	return nets.AllPeers()
}

func IsPeer(id string) bool {
	return nets.IsPeer(id)
}

func GetNodeInfo(id string) (*gen_server.NodeInfo, error) {
	return node.dc.GetNodeInfo(id)
}

func GetNodeStatus(id string, ref interface{}) error {
	return node.dc.GetNodeStatus(id, ref)
}

func Log(obj interface{}) {
	txt, ok := obj.(string)
	if ok {
		node.Logger().Debug(txt)
	} else {
		msg, err := json.Marshal(obj)
		if err != nil {
			fmt.Println("the console data format that cannot be converted!")
		}
		node.Logger().Debug(msg)
	}
}

func LogWarn(msg string) {
	node.Logger().Warn(msg)
}

func LogError(err error) {
	node.Logger().Error(err.Error())
}

func Info() *gen_server.NodeInfo {
	return node.Info()
}

func Self() string {
	return node.Self()
}
