package stl

import (
	"encoding/json"
	"errors"
	"reflect"
	"sync"
)

type Dictionary struct {
	_map map[string]interface{}
	sync.RWMutex
}

func NewDictionary() *Dictionary {
	dict := Dictionary{}
	dict._map = make(map[string]interface{})
	return &dict
}

func NewDictionaryRaw(raw map[string]interface{}) *Dictionary {
	dict := Dictionary{}
	dict._map = raw
	return &dict
}

func (this *Dictionary) Add(key string, value interface{}) error {
	this.Lock()
	defer this.Unlock()

	_, exist := this._map[key]
	if exist {
		return errors.New("Has Contains The Same Key!")
	}
	this._map[key] = value
	return nil
}

func (this *Dictionary) Remove(key string) bool {
	this.Lock()
	defer this.Unlock()

	_, exist := this._map[key]
	if exist {
		delete(this._map, key)
		return true
	}
	return false
}

func (this *Dictionary) Set(key string, value interface{}) {
	this.Lock()
	defer this.Unlock()

	this._map[key] = value
}

func (this *Dictionary) Get(key string) (interface{}, bool) {
	this.RLock()
	defer this.RUnlock()

	v, exist := this._map[key]
	return v, exist
}

func (this *Dictionary) Count() int {
	this.RLock()
	defer this.RUnlock()

	return len(this._map)
}

func (this *Dictionary) ContainsKey(key string) bool {
	this.RLock()
	defer this.RUnlock()

	_, exist := this._map[key]
	return exist
}

func (this *Dictionary) ContainsValue(value interface{}) bool {
	this.RLock()
	defer this.RUnlock()

	for _, v := range this._map {
		if v == value {
			return true
		}
	}
	return false
}

func (this *Dictionary) KeyValuePairs() map[string]interface{} {
	this.RLock()
	defer this.RUnlock()

	return this._map
}

func (this *Dictionary) KeyValueStrings() map[string]string {
	this.RLock()
	defer this.RUnlock()

	temp := make(map[string]string)
	for k, v := range this._map {
		if reflect.TypeOf(v).Name() == "string" {
			temp[k] = v.(string)
		} else {
			data, err := json.Marshal(v)
			if err != nil {
				temp[k] = "nil"
			} else {
				temp[k] = string(data)
			}
		}
	}
	return temp
}

func (this *Dictionary) Clear() {
	this.Lock()
	defer this.Unlock()

	for key, _ := range this._map {
		delete(this._map, key)
	}
}
