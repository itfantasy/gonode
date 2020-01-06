package gonode

import (
	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/core/errs"
	"github.com/itfantasy/gonode/core/goes"
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
func Info() *gen_server.NodeInfo {
	return node.Info()
}

func Self() string {
	return node.Self()
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

func Spawn(fun func([]interface{}), capacity int) int64 {
	return goes.Spawn(fun, capacity)
}
func Spawns(fun func([]interface{}), num int, capacity int) []int64 {
	return goes.Spawns(fun, num, capacity)
}
func Timer(fun func([]interface{}), repeatRate int, args ...interface{}) int64 {
	return goes.Timer(fun, repeatRate, args...)
}
func Go(fun func([]interface{}), args ...interface{}) int64 {
	return goes.Go(fun, args...)
}
func Kill(pid int64) bool {
	return goes.Kill(pid)
}
func Post(pid int64, args ...interface{}) bool {
	return goes.Post(pid, args...)
}
func PostAny(pids []int64, args ...interface{}) bool {
	return goes.PostAny(pids, args...)
}

func Debug(arg0 interface{}, args ...interface{}) {
	node.Logger().Log4Extend(logger.DEBUG, 3, arg0, args...)
}
func LogInfo(arg0 interface{}, args ...interface{}) {
	node.Logger().Log4Extend(logger.INFO, 3, arg0, args...)
}
func LogWarn(arg0 interface{}, args ...interface{}) {
	node.Logger().Log4Extend(logger.WARN, 3, arg0, args...)
}
func LogError(arg0 interface{}, args ...interface{}) {
	node.Logger().Log4Extend(logger.ERROR, 3, arg0, args...)
}
func LogSource(callstack int) string {
	return node.Logger().Source(callstack + 2)
}

func Error(errcode int, errmsg string) error {
	return errs.New(errcode, errmsg)
}
func ErrInfo(err error) (int, string) {
	return errs.Info(err)
}
