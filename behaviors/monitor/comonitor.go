package monitor

import (
	"errors"
	"strconv"

	"github.com/itfantasy/gonode/components"
)

type CoMonitor struct {
	mail   *components.Email
	sendTo string
}

func NewCoMonitor(emailConf string, sendTo string) (*CoMonitor, error) {
	c := new(CoMonitor)
	comp, err := components.NewComponent(emailConf)
	if err != nil {
		return nil, err
	}
	emailComp, ok := comp.(*components.Email)
	if !ok {
		return nil, errors.New("Must Email Type Component!!")
	}
	c.mail = emailComp
	c.sendTo = sendTo
	return c, nil
}

func (c *CoMonitor) GetFrequency() int {
	return 60000
}
func (c *CoMonitor) OnReportError(nodeid string, title string, content string) {
	c.mail.SendTo(c.sendTo, title, "["+nodeid+"]::"+content)
}
func (c *CoMonitor) OnReportCpu(nodeid string, cpu int) {
	if cpu > 80 {
		c.mail.SendTo(c.sendTo, "CpuOverload", "["+nodeid+"]:: The Cpu Has Overloaded!"+strconv.Itoa(cpu)+"%")
	}
}
func (c *CoMonitor) OnReportMemory(nodeid string, mem int) {
	if mem > 3*1024*1024 {
		c.mail.SendTo(c.sendTo, "MemoryOverload", "["+nodeid+"]:: The Memory Has Overloaded!"+strconv.Itoa(mem)+"KB")
	}
}
func (c *CoMonitor) OnCustomEvent(nodeid string, evncode int, content string) {
	c.mail.SendTo(c.sendTo, "["+nodeid+"]::Event #"+strconv.Itoa(evncode), "["+nodeid+"]::"+content)
}
