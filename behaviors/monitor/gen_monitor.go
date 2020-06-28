package monitor

type GenMonitor interface {
	GetFrequency() int
	OnReportError(string, string, string)
	OnReportCpu(string, int)
	OnReportMemory(string, int)
	OnCustomEvent(string, int, string)
}
