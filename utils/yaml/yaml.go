package yaml

import (
	"github.com/itfantasy/gonode/utils/stl"
	_yaml "gopkg.in/yaml.v2"
)

func Encode(obj interface{}) (string, error) {
	b, err := _yaml.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func Decode(str string, ref interface{}) error {
	err := _yaml.Unmarshal([]byte(str), ref)
	return err
}

func ToDict(str string) (*stl.Dictionary, error) {
	var ret map[string]interface{}
	if err := _yaml.Unmarshal([]byte(str), &ret); err != nil {
		return nil, err
	}
	return stl.NewDictionaryRaw(ret), nil
}

func ToList(str string, capacity int) (*stl.List, error) {
	ret := make([]interface{}, 0, capacity)
	if err := _yaml.Unmarshal([]byte(str), &ret); err != nil {
		return nil, err
	}
	return stl.NewListRaw(ret), nil
}
