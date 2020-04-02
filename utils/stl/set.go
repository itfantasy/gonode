package stl

import (
	"sync"
)

type HashSet struct {
	_map map[interface{}]struct{}
	sync.RWMutex
}

func NewHashSet() *HashSet {
	h := HashSet{}
	h._map = make(map[interface{}]struct{})
	return &h
}

func NewHashSetRaw(raw []interface{}) *HashSet {
	h := HashSet{}
	h._map = make(map[interface{}]struct{})
	for _, item := range raw {
		h._map[item] = struct{}{}
	}

	return &h
}

func (h *HashSet) Add(item interface{}) bool {
	h.Lock()
	defer h.Unlock()

	_, exist := h._map[item]
	if exist {
		return false
	}
	h._map[item] = struct{}{}
	return true
}

func (h *HashSet) Remove(item interface{}) bool {
	h.Lock()
	defer h.Unlock()

	_, exist := h._map[item]
	if exist {
		delete(h._map, item)
		return true
	}
	return false
}

func (h *HashSet) Len() int {
	h.RLock()
	defer h.RUnlock()

	return len(h._map)
}

func (h *HashSet) Contains(item interface{}) bool {
	h.RLock()
	defer h.RUnlock()

	_, exist := h._map[item]
	return exist
}

func (h *HashSet) ForEach(fun func(interface{})) {
	h.RLock()
	defer h.RUnlock()

	for k, _ := range h._map {
		fun(k)
	}
}

func (h *HashSet) Items() []interface{} {
	h.RLock()
	defer h.RUnlock()

	items := NewArray(len(h._map))
	for k, _ := range h._map {
		items = append(items, k)
	}
	return items
}

func (h *HashSet) Clear() {
	h.Lock()
	defer h.Unlock()

	for k, _ := range h._map {
		delete(h._map, k)
	}
}

func (h *HashSet) Intersect(another *HashSet) *HashSet {
	clone := h.Clone()
	clone.IntersectWith(another)
	return clone
}

func (h *HashSet) IntersectWith(another *HashSet) {
	h.Lock()
	defer h.Unlock()

	for k, _ := range h._map {
		if !another.Contains(k) {
			delete(h._map, k)
		}
	}
}

func (h *HashSet) Except(another *HashSet) *HashSet {
	clone := h.Clone()
	clone.ExceptWith(another)
	return clone
}

func (h *HashSet) ExceptWith(another *HashSet) {
	h.Lock()
	defer h.Unlock()

	for k, _ := range h._map {
		if another.Contains(k) {
			delete(h._map, k)
		}
	}
}

func (h *HashSet) Union(another *HashSet) *HashSet {
	clone := h.Clone()
	clone.UnionWith(another)
	return clone
}

func (h *HashSet) UnionWith(another *HashSet) {
	another.ForEach(func(item interface{}) {
		if !h.Contains(item) {
			h.Add(item)
		}
	})
}

func (h *HashSet) IsSubset(another *HashSet) bool {
	if h.Len() > another.Len() {
		return false
	}

	h.RLock()
	defer h.RUnlock()

	for k, _ := range h._map {
		if !another.Contains(k) {
			return false
		}
	}
	return true
}

func (h *HashSet) IsSuperset(another *HashSet) bool {
	return another.IsSubset(h)
}

func (h *HashSet) Clone() *HashSet {
	return NewHashSetRaw(h.Items())
}
