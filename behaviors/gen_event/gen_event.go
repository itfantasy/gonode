package gen_event

type GenEventer interface {
	Setup() *EventConf
	OnReportError(string, string)
	OnCpuOverload(int)
	OnMemoryOverload(int)
}
