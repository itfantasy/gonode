package stl

import (
	"encoding/json"
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

func (d *Dictionary) Add(key string, value interface{}) bool {
	d.Lock()
	defer d.Unlock()

	_, exist := d._map[key]
	if exist {
		return false
	}
	d._map[key] = value
	return true
}

func (d *Dictionary) Remove(key string) bool {
	d.Lock()
	defer d.Unlock()

	_, exist := d._map[key]
	if exist {
		delete(d._map, key)
		return true
	}
	return false
}

func (d *Dictionary) Set(key string, value interface{}) {
	d.Lock()
	defer d.Unlock()

	d._map[key] = value
}

func (d *Dictionary) Get(key string) (interface{}, bool) {
	d.RLock()
	defer d.RUnlock()

	v, exist := d._map[key]
	return v, exist
}

func (d *Dictionary) Len() int {
	d.RLock()
	defer d.RUnlock()

	return len(d._map)
}

func (d *Dictionary) ContainsKey(key string) bool {
	d.RLock()
	defer d.RUnlock()

	_, exist := d._map[key]
	return exist
}

func (d *Dictionary) ContainsValue(value interface{}) bool {
	d.RLock()
	defer d.RUnlock()

	for _, v := range d._map {
		if v == value {
			return true
		}
	}
	return false
}

func (d *Dictionary) ForEach(fun func(string, interface{})) {
	d.RLock()
	defer d.RUnlock()

	for k, v := range d._map {
		fun(k, v)
	}
}

func (d *Dictionary) KeyValuePairs() map[string]interface{} {
	d.RLock()
	defer d.RUnlock()

	return d._map
}

func (d *Dictionary) KeyValueStrings() map[string]string {
	d.RLock()
	defer d.RUnlock()

	temp := make(map[string]string)
	for k, v := range d._map {
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

func (d *Dictionary) Clear() {
	d.Lock()
	defer d.Unlock()

	for key, _ := range d._map {
		delete(d._map, key)
	}
}

func (d *Dictionary) ToJson() (string, error) {
	b, err := json.Marshal(d._map)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (d *Dictionary) LoadJson(s string) error {
	return json.Unmarshal([]byte(s), d._map)
}
