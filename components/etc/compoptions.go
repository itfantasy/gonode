package etc

type CompOptions struct {
	options map[string]interface{}
}

func NewCompOptions() *CompOptions {
	componentOptions := new(CompOptions)
	componentOptions.options = make(map[string]interface{})
	return componentOptions
}

func (this *CompOptions) Set(key string, val interface{}) {
	this.options[key] = val
}

func (this *CompOptions) Get(key string) interface{} {
	val, exist := this.options[key]
	if !exist {
		return nil
	}
	return val
}

func (this *CompOptions) GetBool(key string) bool {
	val, exist := this.options[key]
	if !exist {
		return false
	}
	ret, ok := val.(bool)
	if !ok {
		return false
	}
	return ret
}

func (this *CompOptions) GetInt(key string) int {
	val, exist := this.options[key]
	if !exist {
		return 0
	}
	ret, ok := val.(int)
	if !ok {
		return 0
	}
	return ret
}

func (this *CompOptions) GetString(key string) string {
	val, exist := this.options[key]
	if !exist {
		return ""
	}
	ret, ok := val.(string)
	if !ok {
		return ""
	}
	return ret
}

func (this *CompOptions) GetArgs(key string) map[string]interface{} {
	val, exist := this.options[key]
	if !exist {
		return nil
	}
	ret, ok := val.(map[string]interface{})
	if !ok {
		return nil
	}
	return ret
}
