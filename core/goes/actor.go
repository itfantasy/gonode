package goes

import (
	"github.com/itfantasy/gonode/utils/timer"
)

type Actor struct {
	pid       int64
	thefunc   func([]interface{})
	argschan  chan []interface{}
	isKilling bool
}

func newActor(pid int64, fun func([]interface{}), capacity int, repeatRate int, onceArgs []interface{}) *Actor {
	a := new(Actor)
	a.pid = pid
	a.thefunc = fun
	if capacity > 0 {
		a.argschan = make(chan []interface{}, capacity)
	}
	a.isKilling = false

	go func() {
		defer func() {
			if a.argschan != nil {
				close(a.argschan)
			}
			remove(pid)
		}()
		if a.argschan != nil {
			for args := range a.argschan {
				if a.isKilling && args == nil {
					break
				}
				a.do(args)
			}
		} else if repeatRate != 0 {
			for {
				if a.isKilling {
					break
				}
				a.do(onceArgs)
				timer.Sleep(repeatRate)
			}
		} else {
			a.do(onceArgs)
		}
	}()

	return a
}

func (a *Actor) do(args []interface{}) {
	defer AutoRecover(digester)
	a.thefunc(args)
}

func (a *Actor) post(args []interface{}) bool {
	if a.isKilling {
		return false
	}
	a.argschan <- args
	return true
}

func (a *Actor) killing() {
	if !a.isKilling {
		a.isKilling = true
		a.argschan <- nil
	}
}
