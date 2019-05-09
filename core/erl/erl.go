package erl

import (
	"sync"

	"github.com/itfantasy/gonode/core/logger"
	"github.com/itfantasy/gonode/utils/snowflake"
)

var actors sync.Map

func Spawn(fun func([]interface{}), capacity int) int64 {
	pid := snowflake.GenerateRaw()
	actor := newActor(pid, fun, capacity)
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

func Running(pid int64) bool {
	actor, ok := get(pid)
	if !ok {
		return false
	}
	return !actor.isKilling
}

func Waiting(pid int64) int {
	actor, ok := get(pid)
	if !ok {
		return -1
	}
	if actor.isKilling {
		return -1
	}
	return len(actor.argschan)
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

var elogger *logger.Logger

func SetLogger(log *logger.Logger) {
	elogger = log
}
