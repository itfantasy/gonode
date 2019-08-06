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

func Send(id string, msg []byte) {
	node.Send(id, msg)
}

func Close(id string) error {
	return node.Close(id)
}

func AllConnIds() []string {
	return nets.AllConnIds()
}

func Label(id string) string {
	return nets.Label(id)
}

func IsCntConn(id string) bool {
	return nets.IsCntConnId(id)
}

func AllSvcConnIds() []string {
	return nets.AllSvcConnIds()
}

func AllCntConnIds() []string {
	return nets.AllCntConnIds()
}

func Log(obj interface{}) {
	txt, ok := obj.(string)
	if ok {
		node.Logger().Debug(txt)
	} else {
		msg, err := json.Encode(obj)
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
