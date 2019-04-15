package actor

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

func Post(pid uint32, args ...interface{}) bool {
	v, ok := actors.Load(pid)
	if !ok {
		return false
	}
	actor, ok := v.(*Actor)
	if !ok {
		return false
	}
	actor.post(args)
	return true
}
