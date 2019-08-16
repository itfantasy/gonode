package gen_monitor

type GenMonitor interface {
	Setup() *MonitorConf
	OnReportError(string, string, string)
	OnCpuOverload(string, int)
	OnMemoryOverload(string, int)
}
