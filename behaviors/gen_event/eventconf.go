package gen_event

type EventConf struct {
	CpuLimit int
	MemLimit int

	MoniComp string
}

func NewEventConf() *EventConf {
	e := new(EventConf)
	e.CpuLimit = 80
	e.MemLimit = 3 * 1024 * 1024
	return e
}
