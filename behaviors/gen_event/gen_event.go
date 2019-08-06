package gen_event

type GenEvent interface {
	Setup() *EventConf
	OnReportError(string, string, string)
	OnCpuOverload(string, int)
	OnMemoryOverload(string, int)
}
