package stl

func NewMap() map[string]interface{} {
	return make(map[string]interface{})
}

func NewArray(capacity int) []interface{} {
	return make([]interface{}, 0, capacity)
}
