package gonode

import (
	"fmt"
	"runtime"

	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/core/logger"
	"github.com/itfantasy/gonode/nets"
)

// -------------- global ----------------

var node *GoNode = nil

func init() {
	node = &GoNode{}
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

func Info() *gen_server.NodeInfo {
	return node.Info()
}

func Self() string {
	return node.Self()
}

func Log(lv int, arg0 interface{}, args ...interface{}) {
	node.Logger().Log4Extend(lv, 2, arg0, args...)
}
func Debug(arg0 interface{}, args ...interface{}) {
	node.Logger().Log4Extend(logger.DEBUG, 2, arg0, args...)
}
func LogInfo(arg0 interface{}, args ...interface{}) {
	node.Logger().Log4Extend(logger.INFO, 2, arg0, args...)
}
func LogWarn(arg0 interface{}, args ...interface{}) {
	node.Logger().Log4Extend(logger.WARN, 2, arg0, args...)
}
func LogError(arg0 interface{}, args ...interface{}) {
	node.Logger().Log4Extend(logger.ERROR, 2, arg0, args...)
}
func Caller(callstack int) string {
	pc, _, lineno, ok := runtime.Caller(1 + callstack)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineno)
	}
	return src
}
