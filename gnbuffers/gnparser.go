package gnbuffers

import (
	"bytes"
	"encoding/binary"
	"errors"
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

func (this *GnParser) IntArray() ([]int32, error) {
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

func (this *GnParser) Object() (interface{}, error) {
	c, err := this.Byte()
	if err != nil {
		return nil, err
	}
	switch c {
	case 'b':
		return this.Byte()
	case 't':
		return this.Short()
	case 'i':
		return this.Int()
	case 'l':
		return this.Long()
	case 's':
		return this.String()
	case 'H':
		return this.Hash()
	case 'I':
		return this.IntArray()
	default:
		return nil, errors.New("unknow type !!!")
	}
	return nil, errors.New("unknow type !!!")
}
