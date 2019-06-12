package erl

type Actor struct {
	pid       int64
	thefunc   func([]interface{})
	argschan  chan []interface{}
	isKilling bool
}

func newActor(pid int64, fun func([]interface{}), capacity int) *Actor {
	a := new(Actor)
	a.pid = pid
	a.thefunc = fun
	a.argschan = make(chan []interface{}, capacity)
	a.isKilling = false

	go func() {
		defer func() {
			close(a.argschan)
			remove(pid)
		}()

		for args := range a.argschan {
			if a.isKilling && args == nil {
				break
			}
			a.do(args)
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
