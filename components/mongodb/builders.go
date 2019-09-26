package mongodb

import (
	"strconv"
)

const (
	eq     string = "$eq"
	ne            = "$ne"
	gt            = "$gt"
	lt            = "$lt"
	gte           = "$gte"
	lte           = "$lte"
	in            = "$in"
	nin           = "$nin"
	exists        = "$exists"
	regex         = "$regex"
	size          = "$size"
	all           = "$all"
	and           = "$and"
	or            = "$or"
	set           = "$set"
	inc           = "$inc"
	push          = "$push"
	pull          = "$pull"
)

func kv(k string, v interface{}) map[string]interface{} {
	_map := make(map[string]interface{})
	_map[k] = v
	return _map
}

type FilterBuilder struct {
	filters   map[string]interface{}
	curkey    string
	conj      string
	orFilters map[string]interface{}
}

func NewFilterBuilder() *FilterBuilder {
	f := new(FilterBuilder)
	f.filters = make(map[string]interface{})
	return f
}

func (f *FilterBuilder) addFilter(k string, v interface{}) {
	if f.conj == or {
		f.orFilters[k] = v
	} else {
		f.filters[k] = v
	}
}

func (f *FilterBuilder) cloneAndClearOrFilters() map[string]interface{} {
	clone := make(map[string]interface{})
	for k, v := range f.orFilters {
		clone[k] = v
		delete(f.orFilters, k)
	}
	return clone
}

func (f *FilterBuilder) Equal(key string, val interface{}) *FilterBuilder {
	f.addFilter(key, val)
	return f
}

func (f *FilterBuilder) NotEqual(key string, val interface{}) *FilterBuilder {
	f.addFilter(key, kv(ne, val))
	return f
}

func (f *FilterBuilder) Greater(key string, val interface{}) *FilterBuilder {
	f.addFilter(key, kv(gt, val))
	return f
}

func (f *FilterBuilder) GreaterThan(key string, val interface{}) *FilterBuilder {
	f.addFilter(key, kv(gt, val))
	return f
}

func (f *FilterBuilder) LessThan(key string, val interface{}) *FilterBuilder {
	f.addFilter(key, kv(lt, val))
	return f
}

func (f *FilterBuilder) GreaterEqual(key string, val interface{}) *FilterBuilder {
	f.addFilter(key, kv(gte, val))
	return f
}

func (f *FilterBuilder) LessEqual(key string, val interface{}) *FilterBuilder {
	f.addFilter(key, kv(lte, val))
	return f
}

func (f *FilterBuilder) Exists(key string) *FilterBuilder {
	f.addFilter(key, kv(exists, true))
	return f
}

func (f *FilterBuilder) NotExists(key string) *FilterBuilder {
	f.addFilter(key, kv(exists, false))
	return f
}

func (f *FilterBuilder) Regex(key string, val string) *FilterBuilder {
	f.addFilter(key, kv(regex, val))
	return f
}

func (f *FilterBuilder) ArrayIn(key string, val interface{}) *FilterBuilder {
	f.addFilter(key, kv(in, val))
	return f
}

func (f *FilterBuilder) ArrayNotIn(key string, val interface{}) *FilterBuilder {
	f.addFilter(key, kv(nin, val))
	return f
}

func (f *FilterBuilder) ArraySize(key string, num int) *FilterBuilder {
	f.addFilter(key, kv(size, num))
	return f
}

func (f *FilterBuilder) ArrayAll(key string, val interface{}) *FilterBuilder {
	f.addFilter(key, kv(all, val))
	return f
}

func (f *FilterBuilder) ArrayIndex(key string, index int, val interface{}) *FilterBuilder {
	f.addFilter(key+"."+strconv.Itoa(index), val)
	return f
}

func (f *FilterBuilder) And() *FilterBuilder {
	if f.conj == or {
		f.addFilter(or, f.cloneAndClearOrFilters())
	}
	f.conj = and
	return f
}

func (f *FilterBuilder) Or() *FilterBuilder {
	f.conj = or
	return f
}

func (f *FilterBuilder) Serialize() map[string]interface{} {
	if f.conj == or {
		f.addFilter(or, f.cloneAndClearOrFilters())
	}
	f.conj = ""
	return f.filters
}

type OptionBuilder struct {
	options map[string]interface{}
}

func NewOptionBuilder() *OptionBuilder {
	o := new(OptionBuilder)
	o.options = make(map[string]interface{})
	return o
}

func (o *OptionBuilder) addOption(k string, v interface{}) {
	o.options[k] = v
}

func (o *OptionBuilder) Set(key string, val interface{}) *OptionBuilder {
	o.addOption(key, val)
	return o
}

func (o *OptionBuilder) Inc(key string, num int) *OptionBuilder {
	o.addOption(inc, kv(key, num))
	return o
}

func (o *OptionBuilder) Push(key string, val interface{}) *OptionBuilder {
	o.addOption(push, kv(key, val))
	return o
}

func (o *OptionBuilder) Pull(key string, val interface{}) *OptionBuilder {
	o.addOption(pull, kv(key, val))
	return o
}

func (o *OptionBuilder) Serialize() map[string]interface{} {
	return o.options
}
