package erl

import (
	"fmt"
	"runtime/debug"
)

type Actor struct {
	pid       int64
	thefunc   func([]interface{})
	argschan  chan []interface{}
	isKilling bool
}

func newActor(pid int64, fun func([]interface{}), capacity int) *Actor {
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
			this.do(args)
		}
	}()

	return this
}

func (this *Actor) do(args []interface{}) {
	defer func() {
		if err := recover(); err != nil {
			errMsg := "auto recovering..." + fmt.Sprint(err) + "  args:" + fmt.Sprint(args) +
				"\r\n=============== - CallStackInfo - =============== \r\n" + string(debug.Stack())
			if elogger != nil {
				elogger.Error(errMsg)
			} else {
				fmt.Println(errMsg)
			}
		}
	}()
	this.thefunc(args)
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
