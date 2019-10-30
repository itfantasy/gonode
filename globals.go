package gonode

import (
	"errors"
	"strconv"
	"strings"

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
	return node.Logger().Source(callstack + 3)
}
func Error(errcode int, errmsg string) error {
	if errcode == 0 {
		return errors.New(errmsg)
	}
	return errors.New(strconv.Itoa(errcode) + "##" + errmsg)
}
func ErrInfo(err error) (int, string) {
	infos := strings.Split(err.Error(), "##")
	if len(infos) != 2 {
		return 0, err.Error()
	}
	i, err := strconv.Atoi(infos[0])
	if err != nil {
		return 0, err.Error()
	}
	return i, infos[1]
}
