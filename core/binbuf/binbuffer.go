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

func (b *BinBuffer) PushByte(value byte) error {
	return binary.Write(b.bytesBuffer, binary.LittleEndian, value)
}

func (b *BinBuffer) PushBool(value bool) error {
	return binary.Write(b.bytesBuffer, binary.LittleEndian, value)
}

func (b *BinBuffer) PushBytes(value []byte) error {
	return binary.Write(b.bytesBuffer, binary.LittleEndian, value)
}

func (b *BinBuffer) PushShort(value int16) error {
	return binary.Write(b.bytesBuffer, binary.LittleEndian, value)
}

func (b *BinBuffer) PushInt(value int32) error {
	return binary.Write(b.bytesBuffer, binary.LittleEndian, value)
}

func (b *BinBuffer) PushLong(value int64) error {
	return binary.Write(b.bytesBuffer, binary.LittleEndian, value)
}

func (b *BinBuffer) PushString(value string) error {
	buffer := ([]byte)(value)
	if err := b.PushInt(int32(len(buffer))); err != nil { // write the len of the string
		return err
	}
	return binary.Write(b.bytesBuffer, binary.LittleEndian, buffer)
}

func (b *BinBuffer) PushFloat(value float32) error {
	return binary.Write(b.bytesBuffer, binary.LittleEndian, value)
}

func (b *BinBuffer) PushInts(value []int32) error {
	length := len(value)
	if err := b.PushInt(int32(length)); err != nil { // write the len of the []int
		return err
	}
	for _, v := range value {
		if err := b.PushInt(v); err != nil {
			return err
		}
	}
	return nil
}

func (b *BinBuffer) PushArray(value []interface{}) error {
	length := len(value)
	if err := b.PushInt(int32(length)); err != nil { // write the len of the []int
		return err
	}
	for _, v := range value {
		if err := b.PushObject(v); err != nil {
			return err
		}
	}
	return nil
}

func (b *BinBuffer) PushHash(value map[interface{}]interface{}) error {
	length := len(value)
	if err := b.PushInt(int32(length)); err != nil { // write the len of the hash
		return err
	}
	for k, v := range value {
		if err := b.PushObject(k); err != nil {
			return err
		}
		if err := b.PushObject(v); err != nil {
			return err
		}
	}
	return nil
}

func (b *BinBuffer) PushObject(value interface{}) error {
	if value == nil {
		if err := b.PushByte(types.Null); err != nil {
			return err
		}
		if err := b.PushByte(byte(0)); err != nil {
			return err
		}
		return nil
	}
	switch value.(type) {
	case byte:
		if err := b.PushByte(types.Byte); err != nil {
			return err
		}
		if err := b.PushByte(value.(byte)); err != nil {
			return err
		}
	case bool:
		if err := b.PushByte(types.Bool); err != nil {
			return err
		}
		if err := b.PushBool(value.(bool)); err != nil {
			return err
		}
	case int16:
		if err := b.PushByte(types.Short); err != nil {
			return err
		}
		if err := b.PushShort(value.(int16)); err != nil {
			return err
		}
	case int:
		if err := b.PushByte(types.Int); err != nil {
			return err
		}
		if err := b.PushInt(int32(value.(int))); err != nil {
			return err
		}
	case int32:
		if err := b.PushByte(types.Int); err != nil {
			return err
		}
		if err := b.PushInt(value.(int32)); err != nil {
			return err
		}
	case int64:
		if err := b.PushByte(types.Long); err != nil {
			return err
		}
		if err := b.PushLong(value.(int64)); err != nil {
			return err
		}
	case string:
		if err := b.PushByte(byte('s')); err != nil {
			return err
		}
		if err := b.PushString(value.(string)); err != nil {
			return err
		}
	case float32:
		if err := b.PushByte(types.Float); err != nil {
			return err
		}
		if err := b.PushFloat(value.(float32)); err != nil {
			return err
		}
	case []int32:
		if err := b.PushByte(types.Ints); err != nil {
			return err
		}
		if err := b.PushInts(value.([]int32)); err != nil {
			return err
		}
	case []interface{}:
		if err := b.PushByte(types.Array); err != nil {
			return err
		}
		if err := b.PushArray(value.([]interface{})); err != nil {
			return err
		}
	case map[interface{}]interface{}:
		if err := b.PushByte(types.Hash); err != nil {
			return err
		}
		if err := b.PushHash(value.(map[interface{}]interface{})); err != nil {
			return err
		}
	default:
		return errors.New("unsupported type!!")
	}
	return nil
}

func (b *BinBuffer) Bytes() []byte {
	return b.bytesBuffer.Bytes() // has been a slic
}

func (b *BinBuffer) Dispose() {
	b.buffer = nil
}
