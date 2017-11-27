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
	default:
		return nil, errors.New("unknow type !!!")
	}
	return nil, errors.New("unknow type !!!")
}
