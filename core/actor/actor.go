package actor

type Actor struct {
	pid      uint32
	thefunc  func([]interface{})
	argschan chan []interface{}
}

func newActor(pid uint32, fun func([]interface{}), capacity int) *Actor {
	this := new(Actor)
	this.pid = pid
	this.thefunc = fun
	this.argschan = make(chan []interface{}, capacity)

	go func() {
		defer func() {

		}()

		for args := range this.argschan {
			this.thefunc(args)
		}
	}()

	return this
}

func (this *Actor) post(args []interface{}) {
	this.argschan <- args
}
