package binbuf

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
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

func (b *BinBuffer) PushByte(value byte) *BinBuffer {
	if b.err != nil {
		return b
	}
	err := binary.Write(b.bytesBuffer, binary.LittleEndian, value)
	if err != nil {
		b.err = err
		b.errInfo = fmt.Sprintf("PushByte(%v)", value)
	}
	return b
}

func (b *BinBuffer) PushBool(value bool) *BinBuffer {
	if b.err != nil {
		return b
	}
	err := binary.Write(b.bytesBuffer, binary.LittleEndian, value)
	if err != nil {
		b.err = err
		b.errInfo = fmt.Sprintf("PushBool(%v)", value)
	}
	return b
}

func (b *BinBuffer) PushBytes(value []byte) *BinBuffer {
	if b.err != nil {
		return b
	}
	err := binary.Write(b.bytesBuffer, binary.LittleEndian, value)
	if err != nil {
		b.err = err
		b.errInfo = fmt.Sprintf("PushBytes(%v)", value)
	}
	return b
}

func (b *BinBuffer) PushShort(value int16) *BinBuffer {
	if b.err != nil {
		return b
	}
	err := binary.Write(b.bytesBuffer, binary.LittleEndian, value)
	if err != nil {
		b.err = err
		b.errInfo = fmt.Sprintf("PushShort(%v)", value)
	}
	return b
}

func (b *BinBuffer) PushInt(value int32) *BinBuffer {
	if b.err != nil {
		return b
	}
	err := binary.Write(b.bytesBuffer, binary.LittleEndian, value)
	if err != nil {
		b.err = err
		b.errInfo = fmt.Sprintf("PushInt(%v)", value)
	}
	return b
}

func (b *BinBuffer) PushLong(value int64) *BinBuffer {
	if b.err != nil {
		return b
	}
	err := binary.Write(b.bytesBuffer, binary.LittleEndian, value)
	if err != nil {
		b.err = err
		b.errInfo = fmt.Sprintf("PushLong(%v)", value)
	}
	return b
}

func (b *BinBuffer) PushString(value string) *BinBuffer {
	if b.err != nil {
		return b
	}
	buffer := ([]byte)(value)
	b.PushInt(int32(len(buffer))) // write the len of the string
	if b.err != nil {
		return b
	}
	err := binary.Write(b.bytesBuffer, binary.LittleEndian, buffer)
	if err != nil {
		b.err = err
		b.errInfo = fmt.Sprintf("PushString(%v)", value)
	}
	return b
}

func (b *BinBuffer) PushFloat(value float32) *BinBuffer {
	if b.err != nil {
		return b
	}
	err := binary.Write(b.bytesBuffer, binary.LittleEndian, value)
	if err != nil {
		b.err = err
		b.errInfo = fmt.Sprintf("PushFloat(%v)", value)
	}
	return b
}

func (b *BinBuffer) PushInts(value []int32) *BinBuffer {
	if b.err != nil {
		return b
	}
	length := len(value)
	b.PushInt(int32(length)) // write the len of the []int
	if b.err != nil {
		b.errInfo = fmt.Sprintf("PushInts(%v)", value)
		return b
	}
	for _, v := range value {
		b.PushInt(v)
		if b.err != nil {
			b.errInfo = fmt.Sprintf("PushInts(%v)", value)
			return b
		}
	}
	return b
}

func (b *BinBuffer) PushArray(value []interface{}) *BinBuffer {
	if b.err != nil {
		return b
	}
	length := len(value)
	b.PushInt(int32(length)) // write the len of the []int
	if b.err != nil {
		b.errInfo = fmt.Sprintf("PushArray(%v)", value)
		return b
	}
	for _, v := range value {
		b.PushObject(v)
		if b.err != nil {
			b.errInfo = fmt.Sprintf("PushArray(%v)", value)
			return b
		}
	}
	return b
}

func (b *BinBuffer) PushHash(value map[interface{}]interface{}) *BinBuffer {
	if b.err != nil {
		return b
	}
	length := len(value)
	b.PushInt(int32(length)) // write the len of the hash
	if b.err != nil {
		b.errInfo = fmt.Sprintf("PushHash(%v)", value)
		return b
	}
	for k, v := range value {
		b.PushObject(k)
		b.PushObject(v)
		if b.err != nil {
			b.errInfo = fmt.Sprintf("PushHash(%v)", value)
			return b
		}
	}
	return b
}

func (b *BinBuffer) PushObject(value interface{}) *BinBuffer {
	if b.err != nil {
		return b
	}
	if value == nil {
		b.PushByte(Null)
		b.PushByte(byte(0))
		return b
	}
	switch value.(type) {
	case byte:
		b.PushByte(Byte)
		b.PushByte(value.(byte))
	case bool:
		b.PushByte(Bool)
		b.PushBool(value.(bool))
	case int16:
		b.PushByte(Short)
		b.PushShort(value.(int16))
	case int:
		b.PushByte(Int)
		b.PushInt(int32(value.(int)))
	case int32:
		b.PushByte(Int)
		b.PushInt(value.(int32))
	case int64:
		b.PushByte(Long)
		b.PushLong(value.(int64))
	case string:
		b.PushByte(byte('s'))
		b.PushString(value.(string))
	case float32:
		b.PushByte(Float)
		b.PushFloat(value.(float32))
	case []int32:
		b.PushByte(Ints)
		b.PushInts(value.([]int32))
	case []interface{}:
		b.PushByte(Array)
		b.PushArray(value.([]interface{}))
	case map[interface{}]interface{}:
		b.PushByte(Hash)
		b.PushHash(value.(map[interface{}]interface{}))
	default:
		itype := reflect.TypeOf(value)
		ctype, ok := _customBufferExtends[itype]
		if ok {
			b.PushByte(ctype.bSign)
			ctype.serializeFunc(b, value)
		} else {
			b.err = errors.New("unsupported type!!")
		}
	}
	if b.err != nil {
		b.errInfo = fmt.Sprintf("PushObject(%v)", value)
	}
	return b
}

func (b *BinBuffer) Bytes() ([]byte, error) {
	if b.err != nil {
		return nil, b.err
	}
	return b.bytesBuffer.Bytes(), nil // has been a slic
}

func (b *BinBuffer) Dispose() {
	b.buffer = nil
}

func (b *BinBuffer) ErrorInfo() string {
	if b.err != nil {
		return b.errInfo + "|" + b.err.Error()
	}
	return ""
}
