package gnbuffers

import (
	"bytes"
	"encoding/binary"
)

type GnParser struct {
	buffer      []byte
	offset      int
	bytesBuffer *bytes.Buffer
}

func BuildParser(buffer []byte, offset int) *GnParser {
	parser := new(GnParser)
	parser.buffer = buffer
	parser.offset = offset
	parser.bytesBuffer = bytes.NewBuffer(parser.buffer)
	parser.bytesBuffer.Grow(parser.offset)
	return parser
}

func (this *GnParser) Int() (int32, error) {
	var ret int32
	err := binary.Read(this.bytesBuffer, binary.BigEndian, &ret)
	if err != nil {
		return 0, err
	}
	return ret, nil
}

func (this *GnParser) Long() (int64, error) {
	var ret int64
	err := binary.Read(this.bytesBuffer, binary.BigEndian, &ret)
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
	var tempBuffer []byte = make([]byte, length)
	if binary.Read(this.bytesBuffer, binary.BigEndian, &tempBuffer); err != nil {
		return "", err
	}
	return string(tempBuffer), nil

}
