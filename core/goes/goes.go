package goes

import (
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/itfantasy/gonode/utils/snowflake"
)

var actors sync.Map

func Spawn(fun func([]interface{}), capacity int) int64 {
	pid := snowflake.GenerateRaw()
	actor := newActor(pid, fun, capacity, 0, nil)
	actors.Store(pid, actor)
	return pid
}

func Spawns(fun func([]interface{}), num int, capacity int) []int64 {
	pids := make([]int64, 0, num)
	for i := 0; i < num; i++ {
		pids = append(pids, Spawn(fun, capacity))
	}
	return pids
}

func Timer(fun func([]interface{}), repeatRate int, args ...interface{}) int64 {
	pid := snowflake.GenerateRaw()
	actor := newActor(pid, fun, 0, repeatRate, args)
	actors.Store(pid, actor)
	return pid
}

func Go(fun func([]interface{}), args ...interface{}) int64 {
	pid := snowflake.GenerateRaw()
	actor := newActor(pid, fun, 0, 0, args)
	actors.Store(pid, actor)
	return pid
}

func Kill(pid int64) bool {
	actor, ok := get(pid)
	if !ok {
		return false
	}
	actor.killing()
	return true
}

func Post(pid int64, args ...interface{}) bool {
	actor, ok := get(pid)
	if !ok {
		return false
	}
	return actor.post(args)
}

func PostAny(pids []int64, args ...interface{}) bool {
	num := len(pids)
	if num <= 0 {
		return false
	}
	_pid := pids[0]
	_wait := Blocking(_pid)
	for i := 1; i < num; i++ {
		tmpPid := pids[i]
		tmpWait := Blocking(tmpPid)
		if tmpWait < _wait {
			_pid = tmpPid
		}
	}
	return Post(_pid, args...)
}

func Living(pid int64) bool {
	actor, ok := get(pid)
	if !ok {
		return false
	}
	return !actor.isKilling
}

func Blocking(pid int64) int {
	actor, ok := get(pid)
	if !ok {
		return -1
	}
	if actor.isKilling {
		return -1
	}
	return len(actor.argschan)
}

func LivingNum() int {
	return 0
}

func NormalNum() int {
	return 0
}

func BlockingNum() int {
	return 0
}

func get(pid int64) (*Actor, bool) {
	v, ok := actors.Load(pid)
	if !ok {
		return nil, false
	}
	actor, ok := v.(*Actor)
	if !ok {
		return nil, false
	}
	return actor, true
}

func remove(pid int64) {
	actors.Delete(pid)
}

type ErrorDigester interface {
	OnDigestError(interface{})
}

func AutoRecover(e ErrorDigester) {
	if err := recover(); err != nil {
		if e != nil {
			e.OnDigestError(err)
		} else {
			content := "!!! Auto Recovering...  " + fmt.Sprint(err) +
				"\r=============== - CallStackInfo - =============== \r" + string(debug.Stack())
			fmt.Println(content)
		}
	}
}

var digester ErrorDigester

func BindErrorDigester(e ErrorDigester) {
	digester = e
}
