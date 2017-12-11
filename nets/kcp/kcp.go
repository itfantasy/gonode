package kcp

import (
	"github.com/itfantasy/gonode/nets"
	"github.com/xtaci/kcp-go"
)

type KcpNetWorker struct {
	eventListener nets.INetEventListener
}
