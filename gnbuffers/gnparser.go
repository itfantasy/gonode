package gnbuffers

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/itfantasy/gonode/gnbuffers/gntypes"
)

type GnParser struct {
	buffer      []byte
	bytesBuffer *bytes.Buffer
}

func BuildParser(buffer []byte, offset int) *GnParser {
	parser := new(GnParser)
	parser.buffer = buffer
	parser.bytesBuffer = bytes.NewBuffer(parser.buffer)
	parser.bytesBuffer.Grow(offset)
	return parser
}

func (this *GnParser) Byte() (byte, error) {
	var ret byte
	err := binary.Read(this.bytesBuffer, binary.LittleEndian, &ret)
	if err != nil {
		return 0, err
	}
	return ret, nil
}

func (this *GnParser) Bytes() []byte {
	return this.buffer
}

func (this *GnParser) Bool() (bool, error) {
	var ret bool
	err := binary.Read(this.bytesBuffer, binary.LittleEndian, &ret)
	if err != nil {
		return false, err
	}
	return ret, nil
}

func (this *GnParser) Short() (int16, error) {
	var ret int16
	err := binary.Read(this.bytesBuffer, binary.LittleEndian, &ret)
	if err != nil {
		return 0, err
	}
	return ret, nil
}

func (this *GnParser) Int() (int32, error) {
	var ret int32
	err := binary.Read(this.bytesBuffer, binary.LittleEndian, &ret)
	if err != nil {
		return 0, err
	}
	return ret, nil
}

func (this *GnParser) Long() (int64, error) {
	var ret int64
	err := binary.Read(this.bytesBuffer, binary.LittleEndian, &ret)
	if err != nil {
		return 0, err
	}
	return ret, nil
}

func (this *GnParser) String() (string, error) {
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

func (this *GnParser) Float() (float32, error) {
	var ret float32
	err := binary.Read(this.bytesBuffer, binary.LittleEndian, &ret)
	if err != nil {
		return 0, err
	}
	return ret, nil
}

func (this *GnParser) Ints() ([]int32, error) {
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

func (this *GnParser) Array() ([]interface{}, error) {
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

func (this *GnParser) Hash() (map[interface{}]interface{}, error) {
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

func (this *GnParser) Object() (interface{}, error) {
	c, err := this.Byte()
	if err != nil {
		return nil, err
	}
	switch c {
	case gntypes.Byte:
		return this.Byte()
	case gntypes.Bool:
		return this.Bool()
	case gntypes.Short:
		return this.Short()
	case gntypes.Int:
		return this.Int()
	case gntypes.Long:
		return this.Long()
	case gntypes.String:
		return this.String()
	case gntypes.Float:
		return this.Float()
	case gntypes.Ints:
		return this.Ints()
	case gntypes.Array:
		return this.Array()
	case gntypes.Hash:
		return this.Hash()
	default:
		return nil, errors.New("unknow type !!!")
	}
	return nil, errors.New("unknow type !!!")
}
