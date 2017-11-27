package gnbuffers

import (
	"bytes"
	"encoding/binary"
	"errors"
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

func (this *GnBuffer) PushByte(value byte) error {
	return binary.Write(this.bytesBuffer, binary.LittleEndian, value)
}

func (this *GnBuffer) PushShort(value int16) error {
	return binary.Write(this.bytesBuffer, binary.LittleEndian, value)
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

func (this *GnBuffer) PushObject(value interface{}) error {
	switch value.(type) {
	case byte:
		if err := this.PushByte(byte('b')); err != nil {
			return err
		}
		if err := this.PushByte(value.(byte)); err != nil {
			return err
		}
	case int16:
		if err := this.PushByte(byte('t')); err != nil {
			return err
		}
		if err := this.PushShort(value.(int16)); err != nil {
			return err
		}
	case int:
		if err := this.PushByte(byte('i')); err != nil {
			return err
		}
		if err := this.PushInt(int32(value.(int))); err != nil {
			return err
		}
	case int32:
		if err := this.PushByte(byte('i')); err != nil {
			return err
		}
		if err := this.PushInt(value.(int32)); err != nil {
			return err
		}
	case int64:
		if err := this.PushByte(byte('l')); err != nil {
			return err
		}
		if err := this.PushLong(value.(int64)); err != nil {
			return err
		}
	case string:
		if err := this.PushByte(byte('s')); err != nil {
			return err
		}
		if err := this.PushString(value.(string)); err != nil {
			return err
		}
	default:
		return errors.New("unsupported type!!")
	}
	return nil
}

func (this *GnBuffer) Bytes() []byte {
	return this.bytesBuffer.Bytes() // has been a slic
}
