package supervisor

import (
	"errors"

	"github.com/itfantasy/gonode/behaviors/gen_server"
)

const (
	SUPERVISOR string = "supervisor"
	ALLNODES   string = "allnodes"
)

type SuperNode struct {
	super Supervisor
	info  *SuperInfo
}

func NewSuperNode() *SuperNode {
	s := new(SuperNode)
	return s
}

func (s *SuperNode) InitSupervisor(super Supervisor) error {
	s.super = super
	s.info = super.Setup()
	if s.info == nil {
		return errors.New("Can not setup an correct SuperInfo!!")
	}
	return nil
}

func (s *SuperNode) Setup() *gen_server.NodeInfo {
	info := gen_server.NewNodeInfo()
	info.RegDC = s.info.RegDC
	info.NameSpace = s.info.NameSpace
	info.NodeId = SUPERVISOR
	info.IsPub = s.info.IsPub
	info.BackEnds = ALLNODES
	return info
}
func (s *SuperNode) Start() {
	s.super.Start()
}
func (s *SuperNode) OnConn(id string) {
	s.super.OnConn(id)
}
func (s *SuperNode) OnMsg(id string, msg []byte) {
	s.super.OnMsg(id, msg)
}
func (s *SuperNode) OnClose(id string, reason error) {
	s.super.OnClose(id, reason)
}
