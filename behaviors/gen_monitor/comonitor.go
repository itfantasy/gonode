package gen_monitor

import (
	"errors"
	"strconv"

	"github.com/itfantasy/gonode/components"
	"github.com/itfantasy/gonode/components/email"
)

type CoMonitor struct {
	mail   *email.Email
	sendTo string
}

func NewCoMonitor(emailConf string, sendTo string) (*CoMonitor, error) {
	c := new(CoMonitor)
	comp, err := components.NewComponent(emailConf)
	if err != nil {
		return nil, err
	}
	emailComp, ok := comp.(*email.Email)
	if !ok {
		return nil, errors.New("Must Email Type Component!!")
	}
	c.mail = emailComp
	c.sendTo = sendTo
	return c, nil
}

func (c *CoMonitor) Setup() *MonitorConf {
	e := NewMonitorConf()
	return e
}
func (c *CoMonitor) OnReportError(id string, title string, content string) {
	c.mail.SendTo(c.sendTo, title, "["+id+"]::"+content)
}
func (c *CoMonitor) OnCpuOverload(id string, cpu int) {
	c.mail.SendTo(c.sendTo, "CpuOverload", "["+id+"]:: The Cpu Has Overloaded!"+strconv.Itoa(cpu)+"%")
}
func (c *CoMonitor) OnMemoryOverload(id string, mem int) {
	c.mail.SendTo(c.sendTo, "MemoryOverload", "["+id+"]:: The Memory Has Overloaded!"+strconv.Itoa(mem)+"KB")
}
