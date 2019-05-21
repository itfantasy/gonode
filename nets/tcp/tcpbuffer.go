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
	this := new(TcpBuffer)
	this._buffer = buf
	this._offset = 0
	this._count = 0
	this._len = len(this._buffer)
	return this
}

func (this *TcpBuffer) Clear() {
	this._offset = 0
	this._count = 0
	for i := 0; i < this._len; i++ {
		this._buffer[i] = byte(0)
	}
}

func (this *TcpBuffer) Reset() {
	copy(this._buffer, this.Slice())
	this._offset = 0
}

func (this *TcpBuffer) Buffer() []byte {
	return this._buffer[this._offset:]
}

func (this *TcpBuffer) Slice() []byte {
	return this._buffer[this._offset : this._offset+this._count]
}

func (this *TcpBuffer) Count() int {
	return this._count
}

func (this *TcpBuffer) Offset() int {
	return this._offset
}

func (this *TcpBuffer) Capcity() int {
	return this._len - this._offset
}

func (this *TcpBuffer) AddDataLen(count int) {
	this._count += count
}

func (this *TcpBuffer) DeleteData(count int) {
	if this._count >= count {
		this._offset += count
		this._count -= count
	}
}

func (this *TcpBuffer) Dispose() {
	this._buffer = nil
	this._offset = 0
	this._count = 0
	this._len = 0
}
