package sysmonitor

import (
	"strconv"

	"github.com/itfantasy/gonode/behaviors/gen_event"
	"github.com/itfantasy/gonode/core/logger"
	"github.com/itfantasy/gonode/utils/os"
	"github.com/itfantasy/gonode/utils/timer"
)

type SysMonitor struct {
	even    gen_event.GenEventer
	conf    *gen_event.EventConf
	monitor *logger.Logger
}

func NewSysMonitor(id string, even gen_event.GenEventer, monichan string) (*SysMonitor, error) {
	s := new(SysMonitor)
	s.even = even
	s.conf = s.even.Setup()
	if s.conf.MoniComp != "" {
		monitor, err := logger.NewLogger(id, "INFO", monichan, s.conf.MoniComp)
		if err != nil {
			return nil, err
		}
		s.monitor = monitor
	}
	return s, nil
}

func (s *SysMonitor) StartMonitoring() {
	go func() {
		for {
			cpu := int(os.CurCpuPercent())
			mem := int(os.CurMemoryUsage())
			if cpu >= s.conf.CpuLimit {
				s.even.OnCpuOverload(cpu)
			}
			if mem >= s.conf.MemLimit {
				s.even.OnMemoryOverload(mem)
			}
			if s.monitor != nil {
				s.monitor.Info("{\"cpu\":" + strconv.Itoa(cpu) + ",\"mem\":" + strconv.Itoa(mem) + "}")
			}
			timer.Sleep(60000)
		}
	}()
}
