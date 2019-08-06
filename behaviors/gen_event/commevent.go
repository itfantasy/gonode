package gen_event

import (
	"errors"
	"strconv"

	"github.com/itfantasy/gonode/components"
	"github.com/itfantasy/gonode/components/email"
)

type CommEvent struct {
	mail   *email.Email
	sendTo string
}

func NewCommEvent(emailConf string, sendTo string) (*CommEvent, error) {
	c := new(CommEvent)
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

func (c *CommEvent) Setup() *EventConf {
	e := NewEventConf()
	return e
}
func (c *CommEvent) OnReportError(id string, title string, content string) {
	c.mail.SendTo(c.sendTo, title, "["+id+"]::"+content)
}
func (c *CommEvent) OnCpuOverload(id string, cpu int) {
	c.mail.SendTo(c.sendTo, "CpuOverload", "["+id+"]:: The Cpu Has Overloaded!"+strconv.Itoa(cpu)+"%")
}
func (c *CommEvent) OnMemoryOverload(id string, mem int) {
	c.mail.SendTo(c.sendTo, "MemoryOverload", "["+id+"]:: The Memory Has Overloaded!"+strconv.Itoa(mem)+"KB")
}
