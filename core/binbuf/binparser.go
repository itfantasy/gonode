package binbuf

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/itfantasy/gonode/core/binbuf/types"
)

type BinParser struct {
	buffer      []byte
	bytesBuffer *bytes.Buffer
}

func BuildParser(buffer []byte, offset int) *BinParser {
	parser := new(BinParser)
	parser.buffer = buffer
	parser.bytesBuffer = bytes.NewBuffer(parser.buffer)
	parser.bytesBuffer.Grow(offset)
	return parser
}

func (this *BinParser) Byte() (byte, error) {
	var ret byte
	err := binary.Read(this.bytesBuffer, binary.LittleEndian, &ret)
	if err != nil {
		return 0, err
	}
	return ret, nil
}

func (this *BinParser) Bytes() []byte {
	return this.buffer
}

func (this *BinParser) Bool() (bool, error) {
	var ret bool
	err := binary.Read(this.bytesBuffer, binary.LittleEndian, &ret)
	if err != nil {
		return false, err
	}
	return ret, nil
}

func (this *BinParser) Short() (int16, error) {
	var ret int16
	err := binary.Read(this.bytesBuffer, binary.LittleEndian, &ret)
	if err != nil {
		return 0, err
	}
	return ret, nil
}

func (this *BinParser) Int() (int32, error) {
	var ret int32
	err := binary.Read(this.bytesBuffer, binary.LittleEndian, &ret)
	if err != nil {
		return 0, err
	}
	return ret, nil
}

func (this *BinParser) Long() (int64, error) {
	var ret int64
	err := binary.Read(this.bytesBuffer, binary.LittleEndian, &ret)
	if err != nil {
		return 0, err
	}
	return ret, nil
}

func (this *BinParser) String() (string, error) {
	length, err := this.Int() // get the string len
	if err != nil {
		return "", err
	}
	if length > 10240 {
		return "", errors.New("illegal length for a string!!")
	}
	var tempBuffer []byte = make([]byte, length)
	if binary.Read(this.bytesBuffer, binary.LittleEndian, &tempBuffer); err != nil {
		return "", err
	}
	return string(tempBuffer), nil
}

func (this *BinParser) Float() (float32, error) {
	var ret float32
	err := binary.Read(this.bytesBuffer, binary.LittleEndian, &ret)
	if err != nil {
		return 0, err
	}
	return ret, nil
}

func (this *BinParser) Ints() ([]int32, error) {
	length, err := this.Int() // get the []int32 len
	if err != nil {
		return nil, err
	}
	array := make([]int32, 0, length)
	var i int32
	for i = 0; i < length; i++ {
		item, ierr := this.Int()
		if ierr != nil {
			return nil, ierr
		}
		array = append(array, item)
	}
	return array, nil
}

func (this *BinParser) Array() ([]interface{}, error) {
	length, err := this.Int() // get the []int32 len
	if err != nil {
		return nil, err
	}
	array := make([]interface{}, 0, length)
	var i int32
	for i = 0; i < length; i++ {
		item, ierr := this.Object()
		if ierr != nil {
			return nil, ierr
		}
		array = append(array, item)
	}
	return array, nil
}

func (this *BinParser) Hash() (map[interface{}]interface{}, error) {
	length, err := this.Int() // get the hash len
	if err != nil {
		return nil, err
	}
	hash := make(map[interface{}]interface{})
	var i int32
	for i = 0; i < length; i++ {
		k, kerr := this.Object()
		if kerr != nil {
			return nil, kerr
		}
		v, verr := this.Object()
		if verr != nil {
			return nil, verr
		}
		hash[k] = v
	}
	return hash, nil
}

func (this *BinParser) Object() (interface{}, error) {
	c, err := this.Byte()
	if err != nil {
		return nil, err
	}
	switch c {
	case types.Byte:
		return this.Byte()
	case types.Bool:
		return this.Bool()
	case types.Short:
		return this.Short()
	case types.Int:
		return this.Int()
	case types.Long:
		return this.Long()
	case types.String:
		return this.String()
	case types.Float:
		return this.Float()
	case types.Ints:
		return this.Ints()
	case types.Array:
		return this.Array()
	case types.Hash:
		return this.Hash()
	case types.Null:
		if none, err := this.Byte(); err != nil {
			return nil, err
		} else if none != byte(0) {
			return nil, errors.New("unknow type !!!")
		} else {
			return nil, nil
		}
	default:
		return nil, errors.New("unknow type !!!")
	}
	return nil, errors.New("unknow type !!!")
}
