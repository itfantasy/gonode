package tcp

const (
	PCK_MIN_SIZE int   = 6          // |--- header 4bytes ---|--- length 2 bytes ---|--- other datas --- ....
	PCK_HEADER   int32 = 0x2123676f // !#go
)

type TcpBuffer struct {
	_count  int
	_offset int
	_buffer []byte
	_len    int
}

func NewTcpBuffer(buf []byte) *TcpBuffer {
	t := new(TcpBuffer)
	t._buffer = buf
	t._offset = 0
	t._count = 0
	t._len = len(t._buffer)
	return t
}

func (t *TcpBuffer) Clear() {
	t._offset = 0
	t._count = 0
	for i := 0; i < t._len; i++ {
		t._buffer[i] = byte(0)
	}
}

func (t *TcpBuffer) Reset() {
	copy(t._buffer, t.Slice())
	t._offset = 0
}

func (t *TcpBuffer) Buffer() []byte {
	return t._buffer[t._offset:]
}

func (t *TcpBuffer) Slice() []byte {
	return t._buffer[t._offset : t._offset+t._count]
}

func (t *TcpBuffer) Count() int {
	return t._count
}

func (t *TcpBuffer) Offset() int {
	return t._offset
}

func (t *TcpBuffer) Capcity() int {
	return t._len - t._offset
}

func (t *TcpBuffer) AddDataLen(count int) {
	t._count += count
}

func (t *TcpBuffer) DeleteData(count int) {
	if t._count >= count {
		t._offset += count
		t._count -= count
	}
}

func (t *TcpBuffer) Dispose() {
	t._buffer = nil
	t._offset = 0
	t._count = 0
	t._len = 0
}
