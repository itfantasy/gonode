package sysmonitor

import (
	"strconv"

	"github.com/itfantasy/gonode/behaviors/gen_monitor"
	"github.com/itfantasy/gonode/core/logger"
	"github.com/itfantasy/gonode/utils/os"
	"github.com/itfantasy/gonode/utils/timer"
)

type SysMonitoring struct {
	id      string
	monitor gen_monitor.GenMonitor
	conf    *gen_monitor.MonitorConf
	log     *logger.Logger
}

func NewSysMonitoring(id string, monitor gen_monitor.GenMonitor, monichan string) (*SysMonitoring, error) {
	s := new(SysMonitoring)
	s.id = id
	s.monitor = monitor
	s.conf = s.monitor.Setup()
	if s.conf.MoniComp != "" {
		log, err := logger.NewLogger(id, "INFO", monichan, s.conf.MoniComp)
		if err != nil {
			return nil, err
		}
		s.log = log
	}
	return s, nil
}

func (s *SysMonitoring) StartMonitoring() {
	go func() {
		for {
			cpu := int(os.CurCpuPercent())
			mem := int(os.CurMemoryUsage())
			if cpu >= s.conf.CpuLimit {
				s.monitor.OnCpuOverload(s.id, cpu)
			}
			if mem >= s.conf.MemLimit {
				s.monitor.OnMemoryOverload(s.id, mem)
			}
			if s.log != nil {
				s.log.Info("{\"cpu\":" + strconv.Itoa(cpu) + ",\"mem\":" + strconv.Itoa(mem) + "}")
			}
			timer.Sleep(60000)
		}
	}()
}
