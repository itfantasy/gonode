package stl

import (
	"errors"
	"sync"
	//"fmt"
)

type ListInt struct {
	array []int32
	sync.RWMutex
}

func NewListInt(capacity int) *ListInt {
	list := ListInt{}
	list.array = make([]int32, 0, capacity)
	return &list
}

func NewListIntRaw(raw []int32) *ListInt {
	list := ListInt{}
	list.array = raw
	return &list
}

func (this *ListInt) Add(item int32) {
	this.Lock()
	defer this.Unlock()

	this.array = append(this.array, item)
}

func (this *ListInt) Insert(index int, item int32) error {
	this.Lock()
	defer this.Unlock()

	if index > len(this.array) {
		return errors.New("ArgumentOutOfRange")
	}

	temp := make([]int32, 0)
	after := append(temp, this.array[index:]...)
	before := this.array[0:index]
	this.array = append(before, item)
	this.array = append(this.array, after...)
	return nil
}

func (this *ListInt) RemoveAt(index int) error {
	this.Lock()
	defer this.Unlock()

	if index > len(this.array) {
		return errors.New("ArgumentOutOfRange")
	}

	this.array = append(this.array[:index], this.array[index+1:]...)
	return nil
}

func (this *ListInt) Remove(item int32) bool {
	this.Lock()
	defer this.Unlock()

	index := this.IndexOf(item)
	if index < 0 {
		return false
	}
	this.RemoveAt(index)
	return true
}

func (this *ListInt) IndexOf(item int32) int {
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

func (this *ListInt) Contains(item int32) bool {
	this.RLock()
	defer this.RUnlock()

	return this.IndexOf(item) >= 0
}

func (this *ListInt) Count() int {
	this.RLock()
	defer this.RUnlock()

	return len(this.array)
}

func (this *ListInt) Capacity() int {
	this.RLock()
	defer this.RUnlock()

	return cap(this.array)
}

func (this *ListInt) Values() []int32 {
	this.RLock()
	defer this.RUnlock()

	return this.array
}

func (this *ListInt) Get(index int) (int32, error) {
	this.RLock()
	defer this.RUnlock()

	if index > len(this.array) {
		return 0, errors.New("ArgumentOutOfRange")
	}
	return this.array[index], nil
}

func (this *ListInt) Set(index int, item int32) error {
	this.Lock()
	defer this.Unlock()

	if index > len(this.array) {
		return errors.New("ArgumentOutOfRange")
	}
	this.array[index] = item
	return nil
}

func (this *ListInt) Clear() {
	this.Lock()
	defer this.Unlock()

	this.array = this.array[0:0]
}
