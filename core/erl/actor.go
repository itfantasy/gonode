package erl

type Actor struct {
	pid       uint32
	thefunc   func([]interface{})
	argschan  chan []interface{}
	isKilling bool
}

func newActor(pid uint32, fun func([]interface{}), capacity int) *Actor {
	this := new(Actor)
	this.pid = pid
	this.thefunc = fun
	this.argschan = make(chan []interface{}, capacity)
	this.isKilling = false

	go func() {
		defer func() {
			close(this.argschan)
			remove(pid)
		}()

		for args := range this.argschan {
			if this.isKilling && args == nil {
				break
			}
			this.thefunc(args)
		}
	}()

	return this
}

func (this *Actor) post(args []interface{}) bool {
	if this.isKilling {
		return false
	}
	this.argschan <- args
	return true
}

func (this *Actor) killing() {
	if !this.isKilling {
		this.isKilling = true
		this.argschan <- nil
	}
}
