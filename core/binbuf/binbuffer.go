package binbuf

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/itfantasy/gonode/core/binbuf/types"
)

type BinBuffer struct {
	buffer      []byte
	bytesBuffer *bytes.Buffer
	err         error
	errInfo     string
}

func BuildBuffer(capacity int) *BinBuffer {
	buffer := new(BinBuffer)
	buffer.buffer = make([]byte, capacity)
	buffer.bytesBuffer = bytes.NewBuffer(buffer.buffer)
	buffer.bytesBuffer.Reset()
	return buffer
}

func (b *BinBuffer) PushByte(value byte) {
	if b.err != nil {
		return
	}
	err := binary.Write(b.bytesBuffer, binary.LittleEndian, value)
	if err != nil {
		b.err = err
		b.errInfo = fmt.Sprintf("PushByte(%v)", value)
	}
}

func (b *BinBuffer) PushBool(value bool) {
	if b.err != nil {
		return
	}
	err := binary.Write(b.bytesBuffer, binary.LittleEndian, value)
	if err != nil {
		b.err = err
		b.errInfo = fmt.Sprintf("PushBool(%v)", value)
	}
}

func (b *BinBuffer) PushBytes(value []byte) {
	if b.err != nil {
		return
	}
	err := binary.Write(b.bytesBuffer, binary.LittleEndian, value)
	if err != nil {
		b.err = err
		b.errInfo = fmt.Sprintf("PushBytes(%v)", value)
	}
}

func (b *BinBuffer) PushShort(value int16) {
	if b.err != nil {
		return
	}
	err := binary.Write(b.bytesBuffer, binary.LittleEndian, value)
	if err != nil {
		b.err = err
		b.errInfo = fmt.Sprintf("PushShort(%v)", value)
	}
}

func (b *BinBuffer) PushInt(value int32) {
	if b.err != nil {
		return
	}
	err := binary.Write(b.bytesBuffer, binary.LittleEndian, value)
	if err != nil {
		b.err = err
		b.errInfo = fmt.Sprintf("PushInt(%v)", value)
	}
}

func (b *BinBuffer) PushLong(value int64) {
	if b.err != nil {
		return
	}
	err := binary.Write(b.bytesBuffer, binary.LittleEndian, value)
	if err != nil {
		b.err = err
		b.errInfo = fmt.Sprintf("PushLong(%v)", value)
	}
}

func (b *BinBuffer) PushString(value string) {
	if b.err != nil {
		return
	}
	buffer := ([]byte)(value)
	b.PushInt(int32(len(buffer))) // write the len of the string
	if b.err != nil {
		return
	}
	err := binary.Write(b.bytesBuffer, binary.LittleEndian, buffer)
	if err != nil {
		b.err = err
		b.errInfo = fmt.Sprintf("PushString(%v)", value)
	}
}

func (b *BinBuffer) PushFloat(value float32) {
	if b.err != nil {
		return
	}
	err := binary.Write(b.bytesBuffer, binary.LittleEndian, value)
	if err != nil {
		b.err = err
		b.errInfo = fmt.Sprintf("PushFloat(%v)", value)
	}
}

func (b *BinBuffer) PushInts(value []int32) {
	if b.err != nil {
		return
	}
	length := len(value)
	b.PushInt(int32(length)) // write the len of the []int
	if b.err != nil {
		b.errInfo = fmt.Sprintf("PushInts(%v)", value)
		return
	}
	for _, v := range value {
		b.PushInt(v)
		if b.err != nil {
			b.errInfo = fmt.Sprintf("PushInts(%v)", value)
			return
		}
	}
}

func (b *BinBuffer) PushArray(value []interface{}) {
	if b.err != nil {
		return
	}
	length := len(value)
	b.PushInt(int32(length)) // write the len of the []int
	if b.err != nil {
		b.errInfo = fmt.Sprintf("PushArray(%v)", value)
		return
	}
	for _, v := range value {
		b.PushObject(v)
		if b.err != nil {
			b.errInfo = fmt.Sprintf("PushArray(%v)", value)
			return
		}
	}
}

func (b *BinBuffer) PushHash(value map[interface{}]interface{}) {
	if b.err != nil {
		return
	}
	length := len(value)
	b.PushInt(int32(length)) // write the len of the hash
	if b.err != nil {
		b.errInfo = fmt.Sprintf("PushHash(%v)", value)
		return
	}
	for k, v := range value {
		b.PushObject(k)
		b.PushObject(v)
		if b.err != nil {
			b.errInfo = fmt.Sprintf("PushHash(%v)", value)
			return
		}
	}
}

func (b *BinBuffer) PushObject(value interface{}) {
	if b.err != nil {
		return
	}
	if value == nil {
		b.PushByte(types.Null)
		b.PushByte(byte(0))
		return
	}
	switch value.(type) {
	case byte:
		b.PushByte(types.Byte)
		b.PushByte(value.(byte))
	case bool:
		b.PushByte(types.Bool)
		b.PushBool(value.(bool))
	case int16:
		b.PushByte(types.Short)
		b.PushShort(value.(int16))
	case int:
		b.PushByte(types.Int)
		b.PushInt(int32(value.(int)))
	case int32:
		b.PushByte(types.Int)
		b.PushInt(value.(int32))
	case int64:
		b.PushByte(types.Long)
		b.PushLong(value.(int64))
	case string:
		b.PushByte(byte('s'))
		b.PushString(value.(string))
	case float32:
		b.PushByte(types.Float)
		b.PushFloat(value.(float32))
	case []int32:
		b.PushByte(types.Ints)
		b.PushInts(value.([]int32))
	case []interface{}:
		b.PushByte(types.Array)
		b.PushArray(value.([]interface{}))
	case map[interface{}]interface{}:
		b.PushByte(types.Hash)
		b.PushHash(value.(map[interface{}]interface{}))
	default:
		b.err = errors.New("unsupported type!!")
	}
	if b.err != nil {
		b.errInfo = fmt.Sprintf("PushObject(%v)", value)
	}
}

func (b *BinBuffer) Bytes() []byte {
	return b.bytesBuffer.Bytes() // has been a slic
}

func (b *BinBuffer) Dispose() {
	b.buffer = nil
}

func (b *BinBuffer) Error() error {
	return b.err
}

func (b *BinBuffer) ErrorInfo() string {
	if b.err != nil {
		return b.errInfo + "|" + b.err.Error()
	}
	return ""
}
