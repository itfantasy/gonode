package gonode

import (
	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/behaviors/logger"
	"github.com/itfantasy/gonode/behaviors/monitor"
	"github.com/itfantasy/gonode/core/errs"
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
func Info() *gen_server.NodeInfo {
	return node.Info()
}

func Self() string {
	return node.Self()
}

func Bind(behavior interface{}) error {
	return node.Bind(behavior)
}
func BindMonitor(monitor monitor.GenMonitor) {
	node.BindMonitor(monitor)
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
func SendAll(ids []string, msg []byte) []error {
	return node.SendAll(ids, msg)
}
func Close(id string) error {
	return node.Close(id)
}

func Label(id string) string {
	return nets.Label(id)
}
func Nodes(label string) []string {
	return nets.Nodes(label)
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

func Debug(arg0 interface{}, args ...interface{}) {
	node.Logger().Log4Extend(logger.DEBUG, 1, arg0, args...)
}
func LogInfo(arg0 interface{}, args ...interface{}) {
	node.Logger().Log4Extend(logger.INFO, 1, arg0, args...)
}
func LogWarn(arg0 interface{}, args ...interface{}) {
	node.Logger().Log4Extend(logger.WARN, 1, arg0, args...)
}
func LogError(arg0 interface{}, args ...interface{}) {
	node.Logger().Log4Extend(logger.ERROR, 1, arg0, args...)
}
func LogSource(callstack int) string {
	return node.Logger().Source(callstack + 1)
}

func CustomError(errcode int, errmsg string) error {
	return errs.CustomError(errcode, errmsg)
}
func ErrorInfo(err error) (int, string) {
	return errs.ErrorInfo(err)
}
