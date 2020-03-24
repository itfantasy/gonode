package stl

import (
	"sync"
)

type HashTable struct {
	_map map[interface{}]interface{}
	sync.RWMutex
}

func NewHashTable() *HashTable {
	dict := HashTable{}
	dict._map = make(map[interface{}]interface{})
	return &dict
}

func NewHashTableRaw(raw map[interface{}]interface{}) *HashTable {
	dict := HashTable{}
	dict._map = raw
	return &dict
}

func (h *HashTable) Add(key interface{}, value interface{}) bool {
	h.Lock()
	defer h.Unlock()

	_, exist := h._map[key]
	if exist {
		return false
	}
	h._map[key] = value
	return true
}

func (h *HashTable) Remove(key interface{}) bool {
	h.Lock()
	defer h.Unlock()

	_, exist := h._map[key]
	if exist {
		delete(h._map, key)
		return true
	}
	return false
}

func (h *HashTable) Set(key interface{}, value interface{}) {
	h.Lock()
	defer h.Unlock()

	h._map[key] = value
}

func (h *HashTable) Get(key interface{}) (interface{}, bool) {
	h.RLock()
	defer h.RUnlock()

	v, exist := h._map[key]
	return v, exist
}

func (h *HashTable) Len() int {
	h.RLock()
	defer h.RUnlock()

	return len(h._map)
}

func (h *HashTable) ContainsKey(key interface{}) bool {
	h.RLock()
	defer h.RUnlock()

	_, exist := h._map[key]
	return exist
}

func (h *HashTable) ContainsValue(value interface{}) bool {
	h.RLock()
	defer h.RUnlock()

	for _, v := range h._map {
		if v == value {
			return true
		}
	}
	return false
}

func (h *HashTable) ForEach(fun func(interface{}, interface{})) {
	h.RLock()
	defer h.RUnlock()

	for k, v := range h._map {
		fun(k, v)
	}
}

func (h *HashTable) KeyValuePairs() map[interface{}]interface{} {
	h.RLock()
	defer h.RUnlock()

	return h._map
}

func (h *HashTable) Clear() {
	h.Lock()
	defer h.Unlock()

	for key, _ := range h._map {
		delete(h._map, key)
	}
}
