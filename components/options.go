package components

import (
	"strconv"
)

type CompOptions struct {
	options map[string]interface{}
}

func NewCompOptions() *CompOptions {
	componentOptions := new(CompOptions)
	componentOptions.options = make(map[string]interface{})
	return componentOptions
}

func (c *CompOptions) Set(key string, val interface{}) {
	c.options[key] = val
}

func (c *CompOptions) Get(key string) interface{} {
	val, exist := c.options[key]
	if !exist {
		return nil
	}
	return val
}

func (c *CompOptions) GetBool(key string) bool {
	val, exist := c.options[key]
	if !exist {
		return false
	}
	ret, ok := val.(bool)
	if !ok {
		return false
	}
	return ret
}

func (c *CompOptions) GetInt(key string) int {
	val, exist := c.options[key]
	if !exist {
		return 0
	}
	ret, ok := val.(int)
	if !ok {
		strVal, ok := val.(string)
		if !ok {
			return 0
		}
		iVal, err := strconv.Atoi(strVal)
		if err != nil {
			return 0
		}
		return iVal
	}
	return ret
}

func (c *CompOptions) GetString(key string) string {
	val, exist := c.options[key]
	if !exist {
		return ""
	}
	ret, ok := val.(string)
	if !ok {
		return ""
	}
	return ret
}

func (c *CompOptions) GetArgs(key string) map[string]interface{} {
	val, exist := c.options[key]
	if !exist {
		return nil
	}
	ret, ok := val.(map[string]interface{})
	if !ok {
		return nil
	}
	return ret
}
