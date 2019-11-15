package binbuf

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/itfantasy/gonode/core/binbuf/types"
)

type BinParser struct {
	buffer      []byte
	bytesBuffer *bytes.Buffer
	err         error
	errInfo     string
}

func BuildParser(buffer []byte, offset int) *BinParser {
	parser := new(BinParser)
	parser.buffer = buffer
	parser.bytesBuffer = bytes.NewBuffer(parser.buffer)
	parser.bytesBuffer.Grow(offset)
	return parser
}

func (b *BinParser) Byte() byte {
	if b.err != nil {
		return 0
	}
	var ret byte
	err := binary.Read(b.bytesBuffer, binary.LittleEndian, &ret)
	if err != nil {
		b.err = err
		b.errInfo = fmt.Sprintf("Byte()")
		return 0
	}
	return ret
}

func (b *BinParser) Bytes() []byte {
	return b.buffer
}

func (b *BinParser) Bool() bool {
	if b.err != nil {
		return false
	}
	var ret bool
	err := binary.Read(b.bytesBuffer, binary.LittleEndian, &ret)
	if err != nil {
		b.err = err
		b.errInfo = fmt.Sprintf("Bool()")
		return false
	}
	return ret
}

func (b *BinParser) Short() int16 {
	if b.err != nil {
		return 0
	}
	var ret int16
	err := binary.Read(b.bytesBuffer, binary.LittleEndian, &ret)
	if err != nil {
		b.err = err
		b.errInfo = fmt.Sprintf("Short()")
		return 0
	}
	return ret
}

func (b *BinParser) Int() int32 {
	if b.err != nil {
		return 0
	}
	var ret int32
	err := binary.Read(b.bytesBuffer, binary.LittleEndian, &ret)
	if err != nil {
		b.err = err
		b.errInfo = fmt.Sprintf("Int()")
		return 0
	}
	return ret
}

func (b *BinParser) Long() int64 {
	if b.err != nil {
		return 0
	}
	var ret int64
	err := binary.Read(b.bytesBuffer, binary.LittleEndian, &ret)
	if err != nil {
		b.err = err
		b.errInfo = fmt.Sprintf("Long()")
		return 0
	}
	return ret
}

func (b *BinParser) String() string {
	length := b.Int() // get the string len
	if b.err != nil {
		return ""
	}
	if length > 10240 || length < 0 {
		b.err = errors.New("illegal length for a string!!")
		b.errInfo = fmt.Sprintf("String()")
		return ""
	}
	var tempBuffer []byte = make([]byte, length)
	if err := binary.Read(b.bytesBuffer, binary.LittleEndian, &tempBuffer); err != nil {
		b.err = err
		b.errInfo = fmt.Sprintf("String()")
		return ""
	}
	return string(tempBuffer)
}

func (b *BinParser) Float() float32 {
	if b.err != nil {
		return 0
	}
	var ret float32
	err := binary.Read(b.bytesBuffer, binary.LittleEndian, &ret)
	if err != nil {
		b.err = err
		b.errInfo = fmt.Sprintf("Float()")
		return 0
	}
	return ret
}

func (b *BinParser) Ints() []int32 {
	length := b.Int() // get the []int32 len
	if b.err != nil {
		return nil
	}
	array := make([]int32, 0, length)
	var i int32
	for i = 0; i < length; i++ {
		item := b.Int()
		if b.err != nil {
			return nil
		}
		array = append(array, item)
	}
	return array
}

func (b *BinParser) Array() []interface{} {
	length := b.Int() // get the []int32 len
	if b.err != nil {
		return nil
	}
	array := make([]interface{}, 0, length)
	var i int32
	for i = 0; i < length; i++ {
		item := b.Object()
		if b.err != nil {
			return nil
		}
		array = append(array, item)
	}
	return array
}

func (b *BinParser) Hash() map[interface{}]interface{} {
	length := b.Int() // get the hash len
	if b.err != nil {
		return nil
	}
	hash := make(map[interface{}]interface{})
	var i int32
	for i = 0; i < length; i++ {
		k := b.Object()
		v := b.Object()
		if b.err != nil {
			return nil
		}
		hash[k] = v
	}
	return hash
}

func (b *BinParser) Object() interface{} {
	c := b.Byte()
	if b.err != nil {
		return nil
	}
	switch c {
	case types.Byte:
		return b.Byte()
	case types.Bool:
		return b.Bool()
	case types.Short:
		return b.Short()
	case types.Int:
		return b.Int()
	case types.Long:
		return b.Long()
	case types.String:
		return b.String()
	case types.Float:
		return b.Float()
	case types.Ints:
		return b.Ints()
	case types.Array:
		return b.Array()
	case types.Hash:
		return b.Hash()
	case types.Null:
		none := b.Byte()
		if b.err != nil {
			return nil
		} else if none != byte(0) {
			b.err = errors.New("unknow type !!!")
			b.errInfo = fmt.Sprintf("Object()")
			return nil
		} else {
			return nil
		}
	default:
		b.err = errors.New("unknow type !!!")
		b.errInfo = fmt.Sprintf("Object()")
		return nil
	}
	b.err = errors.New("unknow type !!!")
	b.errInfo = fmt.Sprintf("Object()")
	return nil
}

func (b *BinParser) OverFlow() bool {
	return b.bytesBuffer.Len() <= 0 || b.err != nil
}

func (b *BinParser) Error() error {
	return b.err
}

func (b *BinParser) ErrorInfo() string {
	if b.err != nil {
		return b.errInfo + "|" + b.err.Error()
	} else {
		return ""
	}
}
