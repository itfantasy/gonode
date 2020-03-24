package stl

import (
	"encoding/json"
	"errors"
	"sync"
)

type List struct {
	array []interface{}
	sync.RWMutex
}

func NewList(capacity int) *List {
	list := List{}
	list.array = make([]interface{}, 0, capacity)
	return &list
}

func NewListRaw(raw []interface{}) *List {
	list := List{}
	list.array = raw
	return &list
}

func (l *List) Add(item interface{}) {
	l.Lock()
	defer l.Unlock()

	l.array = append(l.array, item)
}

func (l *List) Insert(index int, item interface{}) error {
	l.Lock()
	defer l.Unlock()

	if index > len(l.array) {
		return errors.New("ArgumentOutOfRange")
	}

	temp := make([]interface{}, 0)
	after := append(temp, l.array[index:]...)
	before := l.array[0:index]
	l.array = append(before, item)
	l.array = append(l.array, after...)
	return nil
}

func (l *List) RemoveAt(index int) error {
	l.Lock()
	defer l.Unlock()

	if index > len(l.array) {
		return errors.New("ArgumentOutOfRange")
	}

	l.array = append(l.array[:index], l.array[index+1:]...)
	return nil
}

func (l *List) Remove(item interface{}) bool {
	index := l.IndexOf(item)
	if index < 0 {
		return false
	}
	l.RemoveAt(index)
	return true
}

func (l *List) IndexOf(item interface{}) int {
	l.RLock()
	defer l.RUnlock()

	count := len(l.array)
	for i := 0; i < count; i++ {
		if l.array[i] == item {
			return i
		}
	}
	return -1
}

func (l *List) Contains(item interface{}) bool {
	return l.IndexOf(item) >= 0
}

func (l *List) Len() int {
	l.RLock()
	defer l.RUnlock()

	return len(l.array)
}

func (l *List) Capacity() int {
	l.RLock()
	defer l.RUnlock()

	return cap(l.array)
}

func (l *List) Items() []interface{} {
	l.RLock()
	defer l.RUnlock()

	return l.array
}

func (l *List) Get(index int) (interface{}, error) {
	l.RLock()
	defer l.RUnlock()

	if index >= len(l.array) {
		return nil, errors.New("ArgumentOutOfRange")
	}
	return l.array[index], nil
}

func (l *List) Set(index int, item interface{}) error {
	l.Lock()
	defer l.Unlock()

	if index > len(l.array) {
		return errors.New("ArgumentOutOfRange")
	}
	l.array[index] = item
	return nil
}

func (l *List) ForEach(fun func(interface{})) {
	l.RLock()
	defer l.RUnlock()

	for _, v := range l.array {
		fun(v)
	}
}

func (l *List) Clear() {
	l.Lock()
	defer l.Unlock()

	l.array = l.array[0:0]
}

func (l *List) ToJson() (string, error) {
	b, err := json.Marshal(l.array)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (l *List) LoadJson(s string) error {
	return json.Unmarshal([]byte(s), l.array)
}
