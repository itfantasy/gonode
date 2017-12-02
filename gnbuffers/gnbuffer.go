package gnbuffers

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/itfantasy/gonode/gnbuffers/gntypes"
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

func (this *GnBuffer) PushBytes(value []byte) error {
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

func (this *GnBuffer) PushFloat(value float32) error {
	return binary.Write(this.bytesBuffer, binary.LittleEndian, value)
}

func (this *GnBuffer) PushInts(value []int32) error {
	length := len(value)
	if err := this.PushInt(int32(length)); err != nil { // write the len of the []int
		return err
	}
	for _, v := range value {
		if err := this.PushInt(v); err != nil {
			return err
		}
	}
	return nil
}

func (this *GnBuffer) PushArray(value []interface{}) error {
	length := len(value)
	if err := this.PushInt(int32(length)); err != nil { // write the len of the []int
		return err
	}
	for _, v := range value {
		if err := this.PushObject(v); err != nil {
			return err
		}
	}
	return nil
}

func (this *GnBuffer) PushHash(value map[interface{}]interface{}) error {
	length := len(value)
	if err := this.PushInt(int32(length)); err != nil { // write the len of the hash
		return err
	}
	for k, v := range value {
		if err := this.PushObject(k); err != nil {
			return err
		}
		if err := this.PushObject(v); err != nil {
			return err
		}
	}
	return nil
}

func (this *GnBuffer) PushObject(value interface{}) error {
	switch value.(type) {
	case byte:
		if err := this.PushByte(gntypes.Byte); err != nil {
			return err
		}
		if err := this.PushByte(value.(byte)); err != nil {
			return err
		}
	case int16:
		if err := this.PushByte(gntypes.Short); err != nil {
			return err
		}
		if err := this.PushShort(value.(int16)); err != nil {
			return err
		}
	case int:
		if err := this.PushByte(gntypes.Int); err != nil {
			return err
		}
		if err := this.PushInt(int32(value.(int))); err != nil {
			return err
		}
	case int32:
		if err := this.PushByte(gntypes.Int); err != nil {
			return err
		}
		if err := this.PushInt(value.(int32)); err != nil {
			return err
		}
	case int64:
		if err := this.PushByte(gntypes.Long); err != nil {
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
	case float32:
		if err := this.PushByte(gntypes.Float); err != nil {
			return err
		}
		if err := this.PushFloat(value.(float32)); err != nil {
			return err
		}
	case []int32:
		if err := this.PushByte(gntypes.Ints); err != nil {
			return err
		}
		if err := this.PushInts(value.([]int32)); err != nil {
			return err
		}
	case []interface{}:
		if err := this.PushByte(gntypes.Array); err != nil {
			return err
		}
		if err := this.PushArray(value.([]interface{})); err != nil {
			return err
		}
	case map[interface{}]interface{}:
		if err := this.PushByte(gntypes.Hash); err != nil {
			return err
		}
		if err := this.PushHash(value.(map[interface{}]interface{})); err != nil {
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
