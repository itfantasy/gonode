package supervisor

import (
	"errors"
	"fmt"

	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/components"
	"github.com/itfantasy/gonode/components/rabbitmq"
)

type SuperNode struct {
	id       string
	backEnds string
	super    Supervisor
	info     *SuperInfo
	logComp  *rabbitmq.RabbitMQ
	logChan  string
	moniComp *rabbitmq.RabbitMQ
	moniChan string
}

func NewSuperNode() *SuperNode {
	s := new(SuperNode)
	return s
}

func (s *SuperNode) BindAndInit(superId string, super Supervisor, superBack string, logChan string, moniChan string) error {
	s.id = superId
	s.backEnds = superBack
	s.logChan = logChan
	s.moniChan = moniChan
	s.super = super
	s.info = super.Setup()
	if s.info == nil {
		return errors.New("Can not setup an correct SuperInfo!!")
	}
	comp, err := components.NewComponent(s.info.LogComp)
	if err != nil {
		return err
	}
	rmq, ok := comp.(*rabbitmq.RabbitMQ)
	if !ok {
		return errors.New("The LogComp of SuperInfo MUST be a RabbitMQ Component!!")
	}
	s.logComp = rmq
	s.logComp.BindSubscriber(s)
	comp2, err := components.NewComponent(s.info.MoniComp)
	if err != nil {
		return err
	}
	rmq2, ok := comp2.(*rabbitmq.RabbitMQ)
	if !ok {
		return errors.New("The MoniComp of SuperInfo MUST be a RabbitMQ Component!!")
	}
	s.moniComp = rmq2
	s.moniComp.BindSubscriber(s)
	return nil
}

func (s *SuperNode) Setup() *gen_server.NodeInfo {
	info := gen_server.NewNodeInfo()
	info.Id = s.id
	info.Url = s.info.Url
	info.Pub = false
	info.BackEnds = s.backEnds
	info.LogLevel = "INFO"
	info.LogComp = ""
	info.RegComp = s.info.RegComp
	return info
}
func (s *SuperNode) Start() {
	go s.logComp.Subscribe(s.logChan)
	go s.moniComp.Subscribe(s.moniChan)
	s.super.Start()
}
func (s *SuperNode) OnConn(id string) {
	s.super.OnConn(id)
}
func (s *SuperNode) OnMsg(id string, msg []byte) {
	// nothing to do ...
}
func (s *SuperNode) OnClose(id string, reason error) {
	s.super.OnClose(id, reason)
}

func (s *SuperNode) OnSubscribe(channel string) {
	fmt.Println("SuperNode.OnSubscribe::" + channel)
}
func (s *SuperNode) OnSubMessage(channel string, msg string) {
	fmt.Println("SuperNode.OnSubMessage::" + channel + ":" + msg)
}
func (s *SuperNode) OnSubError(channel string, err error) {
	fmt.Println("SuperNode.OnSubError::" + channel + ":" + err.Error())
}
