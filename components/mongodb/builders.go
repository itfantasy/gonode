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

	set  = "$set"
	inc  = "$inc"
	push = "$push"
	pull = "$pull"

	group    = "$group"
	sum      = "$sum"
	avg      = "$avg"
	min      = "$min"
	max      = "$max"
	addToSet = "$addToSet"
	first    = "$first"
	last     = "$last"
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

func NewFilter() *FilterBuilder {
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

func NewOption() *OptionBuilder {
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

type GroupBuilder struct {
	groups map[interface{}]map[string]interface{}
}

func NewGroup() *GroupBuilder {
	g := new(GroupBuilder)
	g.groups = make(map[interface{}]map[string]interface{})
	return g
}

func (g *GroupBuilder) groupby(by interface{}) interface{} {
	strby, ok := by.(string)
	if ok {
		return "$" + strby
	}
	return by
}

func (g *GroupBuilder) addgroup(opt string, k interface{}, by interface{}, retfield string) *GroupBuilder {
	by := g.groupby(by)
	_, exist := g.groups[by]
	if !exist {
		g.groups[by] = kv("_id", by)
	}
	_map, _ := g.groups[by]
	_map["retfield"] = kv(opt, k)
	return g
}

func (g *GroupBuilder) CountBy(by interface{}) *GroupBuilder {
	return g.addgroup(sum, 1, by)
}

func (g *GroupBuilder) SumBy(key interface{}, by interface{}) *GroupBuilder {
	return g.addgroup(sum, key, by)
}

func (g *GroupBuilder) AvgBy(key interface{}, by interface{}) *GroupBuilder {
	return g.addgroup(avg, key, by)
}

func (g *GroupBuilder) MinBy(key interface{}, by interface{}) *GroupBuilder {
	return g.addgroup(min, key, by)
}

func (g *GroupBuilder) MaxBy(key interface{}, by interface{}) *GroupBuilder {
	return g.addgroup(min, key, by)
}

func (g *GroupBuilder) PushBy(key interface{}, by interface{}) *GroupBuilder {
	return g.addgroup(push, key, by)
}

func (g *GroupBuilder) AddToSetBy(key interface{}, by interface{}) *GroupBuilder {
	return g.addgroup(addToSet, key, by)
}

func (g *GroupBuilder) FirstBy(key interface{}, by interface{}) *GroupBuilder {
	return g.addgroup(first, key, by)
}

func (g *GroupBuilder) LastBy(key interface{}, by interface{}) *GroupBuilder {
	return g.addgroup(last, key, by)
}

func (g *GroupBuilder) Serialize() []interface{} {
	array := make([]interface{}, 0, len(g.groups))
	for _, v := range g.groups {
		array = append(array, v)
	}
	g.groups = nil
	return array
}
