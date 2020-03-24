package actors

import (
	"sync"

	"github.com/itfantasy/gonode/utils/snowflake"
)

var _executors sync.Map

func Spawn(capacity int) *Executor {
	pid := snowflake.GenerateRaw()
	executor := newexecutor(pid, capacity)
	_executors.Store(pid, executor)
	return executor
}

func Kill(pid int64) bool {
	executor, ok := Get(pid)
	if !ok {
		return false
	}
	executor.Dispose()
	return true
}

func Post(pid int64, act func()) bool {
	executor, ok := Get(pid)
	if !ok {
		return false
	}
	return executor.Execute(act)
}

func PostAll(pids []int64, act func()) bool {
	for _, pid := range pids {
		ret := Post(pid, act)
		if !ret {
			return false
		}
	}
	return true
}

func PostAny(pids []int64, act func()) bool {
	num := len(pids)
	if num <= 0 {
		return false
	}
	_pid := pids[0]
	_wait := 999999
	for i := 0; i < num; i++ {
		tmpPid := pids[i]
		tmpWait := TaskLen(tmpPid)
		if tmpWait > 0 && tmpWait < _wait {
			_pid = tmpPid
		}
	}
	return Post(_pid, act)
}

func Living(pid int64) bool {
	executor, ok := Get(pid)
	if !ok {
		return false
	}
	return executor.Living()
}

func TaskLen(pid int64) int {
	executor, ok := Get(pid)
	if !ok {
		return -1
	}
	if executor.killing {
		return -1
	}
	return executor.TaskLen()
}

func Get(pid int64) (*Executor, bool) {
	v, ok := _executors.Load(pid)
	if !ok {
		return nil, false
	}
	executor, ok := v.(*Executor)
	if !ok {
		return nil, false
	}
	return executor, true
}

func remove(pid int64) {
	_executors.Delete(pid)
}
