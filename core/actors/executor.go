package actors

import (
	"github.com/itfantasy/gonode/core/errs"
)

type Executor struct {
	pid     int64
	actchan chan func()
	killing bool
	binder  interface{}
}

func newexecutor(pid int64, capacity int) *Executor {
	e := new(Executor)
	e.pid = pid
	e.actchan = make(chan func(), capacity)
	e.killing = false

	go func() {
		defer func() {
			if e.actchan != nil {
				close(e.actchan)
				e.actchan = nil
			}
			remove(pid)
		}()
		if e.actchan != nil {
			for act := range e.actchan {
				if e.killing && act == nil {
					break
				}
				if act != nil {
					e.doact(act)
				}
			}
		}
	}()

	return e
}

func (e *Executor) Pid() int64 {
	return e.pid
}

func (e *Executor) Execute(act func()) bool {
	if e.killing {
		return false
	}
	e.actchan <- act
	return true
}

func (e *Executor) doact(act func()) {
	defer errs.AutoRecover()
	act()
}

func (e *Executor) Dispose() {
	if !e.killing {
		e.killing = true
		e.actchan <- nil
	}
}

func (e *Executor) Living() bool {
	return !e.killing
}

func (e *Executor) TaskLen() int {
	return len(e.actchan)
}

func (e *Executor) Binder() interface{} {
	return e.binder
}

func (e *Executor) SetBinder(obj interface{}) {
	e.binder = obj
}
