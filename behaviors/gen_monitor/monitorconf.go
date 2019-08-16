package gen_monitor

type MonitorConf struct {
	CpuLimit int
	MemLimit int

	MoniComp string
}

func NewMonitorConf() *MonitorConf {
	e := new(MonitorConf)
	e.CpuLimit = 80
	e.MemLimit = 3 * 1024 * 1024
	return e
}
