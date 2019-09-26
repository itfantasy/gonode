package json

import (
	"fmt"

	"github.com/itfantasy/gonode/utils/stl"
	"github.com/json-iterator/go"
)

func Marshal(obj interface{}) (string, error) {
	b, err := jsoniter.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func Unmarshal(str string, ref interface{}) error {
	err := jsoniter.Unmarshal([]byte(str), ref)
	return err
}

func ToDict(str string) (*stl.Dictionary, error) {
	var ret map[string]interface{}
	if err := jsoniter.Unmarshal([]byte(str), &ret); err != nil {
		return nil, err
	}
	return stl.NewDictionaryRaw(ret), nil
}

func ToList(str string, capacity int) (*stl.List, error) {
	ret := make([]interface{}, 0, capacity)
	if err := jsoniter.Unmarshal([]byte(str), &ret); err != nil {
		return nil, err
	}
	return stl.NewListRaw(ret), nil
}

func Println(obj interface{}) {
	ret, err := Marshal(obj)
	if err != nil {
		fmt.Println("ERROR Data...")
	} else {
		fmt.Println(ret)
	}
}
