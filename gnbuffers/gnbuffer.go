package gnbuffers

import (
	"bytes"
	"encoding/binary"
)

type GnBuffer struct {
	buffer      []byte
	bytesBuffer *bytes.Buffer
}

func BuildBuffer(capacity int) (*GnBuffer, error) {
	if capacity > 10240 {
		return nil, errors.New("illegal length for the buffer!!")
	}
	buffer := new(GnBuffer)
	buffer.buffer = make([]byte, capacity)
	buffer.bytesBuffer = bytes.NewBuffer(buffer.buffer)
	buffer.bytesBuffer.Reset()
	return buffer, nil
}

func (this *GnBuffer) PushInt(value int32) error {
	return binary.Write(this.bytesBuffer, binary.LittleEndian, value)
}

func (this *GnBuffer) PushLong(value int64) error {
	return binary.Write(this.bytesBuffer, binary.LittleEndian, value)
}

func (this *GnBuffer) PushString(value string) error {
	buffer := ([]byte)(value)
	if err := this.PushInt(int32(len(buffer))); err != nil { // write the len of the string
		return err
	}
	return binary.Write(this.bytesBuffer, binary.LittleEndian, buffer)
}

func (this *GnBuffer) Bytes() []byte {
	return this.bytesBuffer.Bytes() // has been a slic
}
