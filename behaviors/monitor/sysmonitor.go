package monitor

import (
	"github.com/itfantasy/gonode/utils/os"
	"github.com/itfantasy/gonode/utils/timer"
)

type SysMonitoring struct {
	id      string
	monitor GenMonitor
}

func NewSysMonitoring(id string, monitor GenMonitor) (*SysMonitoring, error) {
	s := new(SysMonitoring)
	s.id = id
	s.monitor = monitor
	return s, nil
}

func (s *SysMonitoring) StartMonitoring() {
	go func() {
		for {
			cpu := int(os.CurCpuPercent())
			mem := int(os.CurMemoryUsage())
			s.monitor.OnReportCpu(s.id, cpu)
			s.monitor.OnReportMemory(s.id, mem)
			timer.Sleep(s.monitor.GetFrequency())
		}
	}()
}
