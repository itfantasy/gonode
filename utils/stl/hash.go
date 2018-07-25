package stl

import (
	"errors"
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

func (this *HashTable) Add(key interface{}, value interface{}) error {
	this.Lock()
	defer this.Unlock()

	_, exist := this._map[key]
	if exist {
		return errors.New("Has Contains The Same Key!")
	}
	this._map[key] = value
	return nil
}

func (this *HashTable) Remove(key interface{}) error {
	this.Lock()
	defer this.Unlock()

	_, exist := this._map[key]
	if exist {
		delete(this._map, key)
		return nil
	}
	return errors.New("Do Not Has The Key!")
}

func (this *HashTable) Set(key interface{}, value interface{}) {
	this.Lock()
	defer this.Unlock()

	this._map[key] = value
}

func (this *HashTable) Get(key interface{}) (interface{}, bool) {
	this.RLock()
	defer this.RUnlock()

	v, exist := this._map[key]
	return v, exist
}

func (this *HashTable) Count() int {
	this.RLock()
	defer this.RUnlock()

	return len(this._map)
}

func (this *HashTable) ContainsKey(key interface{}) bool {
	this.RLock()
	defer this.RUnlock()

	_, exist := this._map[key]
	return exist
}

func (this *HashTable) ContainsValue(value interface{}) bool {
	this.RLock()
	defer this.RUnlock()

	for _, v := range this._map {
		if v == value {
			return true
		}
	}
	return false
}

func (this *HashTable) KeyValuePairs() map[interface{}]interface{} {
	this.RLock()
	defer this.RUnlock()

	return this._map
}

func (this *HashTable) Clear() {
	this.Lock()
	defer this.Unlock()

	for key, _ := range this._map {
		delete(this._map, key)
	}
}
