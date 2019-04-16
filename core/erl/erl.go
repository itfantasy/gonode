package erl

import (
	"hash/crc32"
	"sync"

	"github.com/itfantasy/gonode/utils/crypt"
)

var actors sync.Map

func Spawn(fun func([]interface{}), capacity int) uint32 {
	pid := crc32.ChecksumIEEE([]byte(crypt.Guid()))
	actor := newActor(pid, fun, capacity)
	actors.Store(pid, actor)
	return pid
}

func Kill(pid uint32) bool {
	actor, ok := get(pid)
	if !ok {
		return false
	}
	actor.killing()
	return true
}

func Post(pid uint32, args ...interface{}) bool {
	actor, ok := get(pid)
	if !ok {
		return false
	}
	return actor.post(args)
}

func Running(pid uint32) bool {
	actor, ok := get(pid)
	if !ok {
		return false
	}
	return !actor.isKilling
}

func Waiting(pid uint32) int {
	actor, ok := get(pid)
	if !ok {
		return -1
	}
	if actor.isKilling {
		return -1
	}
	return len(actor.argschan)
}

func get(pid uint32) (*Actor, bool) {
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

func remove(pid uint32) {
	actors.Delete(pid)
}
