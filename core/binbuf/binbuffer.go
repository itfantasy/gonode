package binbuf

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/itfantasy/gonode/core/binbuf/types"
)

type BinBuffer struct {
	buffer      []byte
	bytesBuffer *bytes.Buffer
}

func BuildBuffer(capacity int) (*BinBuffer, error) {
	if capacity > 10240 {
		return nil, errors.New("illegal length for the buffer!!")
	}
	buffer := new(BinBuffer)
	buffer.buffer = make([]byte, capacity)
	buffer.bytesBuffer = bytes.NewBuffer(buffer.buffer)
	buffer.bytesBuffer.Reset()
	return buffer, nil
}

func (this *BinBuffer) PushByte(value byte) error {
	return binary.Write(this.bytesBuffer, binary.LittleEndian, value)
}

func (this *BinBuffer) PushBool(value bool) error {
	return binary.Write(this.bytesBuffer, binary.LittleEndian, value)
}

func (this *BinBuffer) PushBytes(value []byte) error {
	return binary.Write(this.bytesBuffer, binary.LittleEndian, value)
}

func (this *BinBuffer) PushShort(value int16) error {
	return binary.Write(this.bytesBuffer, binary.LittleEndian, value)
}

func (this *BinBuffer) PushInt(value int32) error {
	return binary.Write(this.bytesBuffer, binary.LittleEndian, value)
}

func (this *BinBuffer) PushLong(value int64) error {
	return binary.Write(this.bytesBuffer, binary.LittleEndian, value)
}

func (this *BinBuffer) PushString(value string) error {
	buffer := ([]byte)(value)
	if err := this.PushInt(int32(len(buffer))); err != nil { // write the len of the string
		return err
	}
	return binary.Write(this.bytesBuffer, binary.LittleEndian, buffer)
}

func (this *BinBuffer) PushFloat(value float32) error {
	return binary.Write(this.bytesBuffer, binary.LittleEndian, value)
}

func (this *BinBuffer) PushInts(value []int32) error {
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

func (this *BinBuffer) PushArray(value []interface{}) error {
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

func (this *BinBuffer) PushHash(value map[interface{}]interface{}) error {
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

func (this *BinBuffer) PushObject(value interface{}) error {
	if value == nil {
		if err := this.PushByte(types.Null); err != nil {
			return err
		}
		if err := this.PushByte(byte(0)); err != nil {
			return err
		}
		return nil
	}
	switch value.(type) {
	case byte:
		if err := this.PushByte(types.Byte); err != nil {
			return err
		}
		if err := this.PushByte(value.(byte)); err != nil {
			return err
		}
	case bool:
		if err := this.PushByte(types.Bool); err != nil {
			return err
		}
		if err := this.PushBool(value.(bool)); err != nil {
			return err
		}
	case int16:
		if err := this.PushByte(types.Short); err != nil {
			return err
		}
		if err := this.PushShort(value.(int16)); err != nil {
			return err
		}
	case int:
		if err := this.PushByte(types.Int); err != nil {
			return err
		}
		if err := this.PushInt(int32(value.(int))); err != nil {
			return err
		}
	case int32:
		if err := this.PushByte(types.Int); err != nil {
			return err
		}
		if err := this.PushInt(value.(int32)); err != nil {
			return err
		}
	case int64:
		if err := this.PushByte(types.Long); err != nil {
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
		if err := this.PushByte(types.Float); err != nil {
			return err
		}
		if err := this.PushFloat(value.(float32)); err != nil {
			return err
		}
	case []int32:
		if err := this.PushByte(types.Ints); err != nil {
			return err
		}
		if err := this.PushInts(value.([]int32)); err != nil {
			return err
		}
	case []interface{}:
		if err := this.PushByte(types.Array); err != nil {
			return err
		}
		if err := this.PushArray(value.([]interface{})); err != nil {
			return err
		}
	case map[interface{}]interface{}:
		if err := this.PushByte(types.Hash); err != nil {
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

func (this *BinBuffer) Bytes() []byte {
	return this.bytesBuffer.Bytes() // has been a slic
}

func (this *BinBuffer) Dispose() {
	this.buffer = nil
}
