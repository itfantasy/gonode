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

func (b *BinParser) Byte() (byte, error) {
	var ret byte
	err := binary.Read(b.bytesBuffer, binary.LittleEndian, &ret)
	if err != nil {
		return 0, err
	}
	return ret, nil
}

func (b *BinParser) Bytes() []byte {
	return b.buffer
}

func (b *BinParser) Bool() (bool, error) {
	var ret bool
	err := binary.Read(b.bytesBuffer, binary.LittleEndian, &ret)
	if err != nil {
		return false, err
	}
	return ret, nil
}

func (b *BinParser) Short() (int16, error) {
	var ret int16
	err := binary.Read(b.bytesBuffer, binary.LittleEndian, &ret)
	if err != nil {
		return 0, err
	}
	return ret, nil
}

func (b *BinParser) Int() (int32, error) {
	var ret int32
	err := binary.Read(b.bytesBuffer, binary.LittleEndian, &ret)
	if err != nil {
		return 0, err
	}
	return ret, nil
}

func (b *BinParser) Long() (int64, error) {
	var ret int64
	err := binary.Read(b.bytesBuffer, binary.LittleEndian, &ret)
	if err != nil {
		return 0, err
	}
	return ret, nil
}

func (b *BinParser) String() (string, error) {
	length, err := b.Int() // get the string len
	if err != nil {
		return "", err
	}
	if length > 10240 {
		return "", errors.New("illegal length for a string!!")
	}
	var tempBuffer []byte = make([]byte, length)
	if binary.Read(b.bytesBuffer, binary.LittleEndian, &tempBuffer); err != nil {
		return "", err
	}
	return string(tempBuffer), nil
}

func (b *BinParser) Float() (float32, error) {
	var ret float32
	err := binary.Read(b.bytesBuffer, binary.LittleEndian, &ret)
	if err != nil {
		return 0, err
	}
	return ret, nil
}

func (b *BinParser) Ints() ([]int32, error) {
	length, err := b.Int() // get the []int32 len
	if err != nil {
		return nil, err
	}
	array := make([]int32, 0, length)
	var i int32
	for i = 0; i < length; i++ {
		item, ierr := b.Int()
		if ierr != nil {
			return nil, ierr
		}
		array = append(array, item)
	}
	return array, nil
}

func (b *BinParser) Array() ([]interface{}, error) {
	length, err := b.Int() // get the []int32 len
	if err != nil {
		return nil, err
	}
	array := make([]interface{}, 0, length)
	var i int32
	for i = 0; i < length; i++ {
		item, ierr := b.Object()
		if ierr != nil {
			return nil, ierr
		}
		array = append(array, item)
	}
	return array, nil
}

func (b *BinParser) Hash() (map[interface{}]interface{}, error) {
	length, err := b.Int() // get the hash len
	if err != nil {
		return nil, err
	}
	hash := make(map[interface{}]interface{})
	var i int32
	for i = 0; i < length; i++ {
		k, kerr := b.Object()
		if kerr != nil {
			return nil, kerr
		}
		v, verr := b.Object()
		if verr != nil {
			return nil, verr
		}
		hash[k] = v
	}
	return hash, nil
}

func (b *BinParser) Object() (interface{}, error) {
	c, err := b.Byte()
	if err != nil {
		return nil, err
	}
	switch c {
	case types.Byte:
		return b.Byte()
	case types.Bool:
		return b.Bool()
	case types.Short:
		return b.Short()
	case types.Int:
		return b.Int()
	case types.Long:
		return b.Long()
	case types.String:
		return b.String()
	case types.Float:
		return b.Float()
	case types.Ints:
		return b.Ints()
	case types.Array:
		return b.Array()
	case types.Hash:
		return b.Hash()
	case types.Null:
		if none, err := b.Byte(); err != nil {
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
