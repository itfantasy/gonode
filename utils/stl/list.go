package stl

import (
	"errors"
	"sync"
	//"fmt"
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

func (this *List) Add(item interface{}) {
	this.Lock()
	defer this.Unlock()

	this.array = append(this.array, item)
}

func (this *List) Insert(index int, item interface{}) error {
	this.Lock()
	defer this.Unlock()

	if index > len(this.array) {
		return errors.New("ArgumentOutOfRange")
	}

	temp := make([]interface{}, 0)
	after := append(temp, this.array[index:]...)
	before := this.array[0:index]
	this.array = append(before, item)
	this.array = append(this.array, after...)
	return nil
}

func (this *List) RemoveAt(index int) error {
	this.Lock()
	defer this.Unlock()

	if index > len(this.array) {
		return errors.New("ArgumentOutOfRange")
	}

	this.array = append(this.array[:index], this.array[index+1:]...)
	return nil
}

func (this *List) Remove(item interface{}) bool {
	this.Lock()
	defer this.Unlock()

	index := this.IndexOf(item)
	if index < 0 {
		return false
	}
	this.RemoveAt(index)
	return true
}

func (this *List) IndexOf(item interface{}) int {
	this.RLock()
	defer this.RUnlock()

	count := len(this.array)
	for i := 0; i < count; i++ {
		if this.array[i] == item {
			return i
		}
	}
	return -1
}

func (this *List) Contains(item interface{}) bool {
	this.RLock()
	defer this.RUnlock()

	return this.IndexOf(item) >= 0
}

func (this *List) Count() int {
	this.RLock()
	defer this.RUnlock()

	return len(this.array)
}

func (this *List) Capacity() int {
	this.RLock()
	defer this.RUnlock()

	return cap(this.array)
}

func (this *List) Values() []interface{} {
	this.RLock()
	defer this.RUnlock()

	return this.array
}

func (this *List) Get(index int) (interface{}, error) {
	this.RLock()
	defer this.RUnlock()

	if index > len(this.array) {
		return nil, errors.New("ArgumentOutOfRange")
	}
	return this.array[index], nil
}

func (this *List) Set(index int, item interface{}) error {
	this.Lock()
	defer this.Unlock()

	if index > len(this.array) {
		return errors.New("ArgumentOutOfRange")
	}
	this.array[index] = item
	return nil
}

func (this *List) Clear() {
	this.Lock()
	defer this.Unlock()

	this.array = this.array[0:0]
}
