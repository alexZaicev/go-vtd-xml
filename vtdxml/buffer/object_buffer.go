package buffer

import (
	"github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"
)

const (
	DefaultObjectPageSize = 1024
)

type FastObjectBufferOption func(*FastObjectBuffer)

type ObjectBuffer interface {
	Buffer
	ObjectAt(index int) (interface{}, error)
	ModifyEntry(index int, value interface{}) error
	Append(value interface{}) error
}

type FastObjectBuffer struct {
	buffer   [][]interface{}
	capacity int
	size     int
	exp      int
	r        int
	pageSize int
}

func WithFastObjectBufferPageSize(size int) FastObjectBufferOption {
	return func(b *FastObjectBuffer) {
		b.pageSize = 1 << size
		b.exp = size
		b.r = b.pageSize - 1
	}
}

func NewFastObjectBuffer(opts ...FastObjectBufferOption) (*FastObjectBuffer, error) {
	b := &FastObjectBuffer{
		capacity: 0,
		size:     0,
		pageSize: DefaultObjectPageSize,
		exp:      10,
		r:        DefaultObjectPageSize - 1,
		buffer:   make([][]interface{}, 0, DefaultObjectPageSize),
	}

	for _, opt := range opts {
		opt(b)
	}
	if b.pageSize == 0 || b.r < 0 {
		return nil, erroring.NewInvalidArgumentError("size", erroring.InvalidBufferPageSize, nil)
	}

	return b, nil
}

func (b *FastObjectBuffer) ObjectAt(index int) (interface{}, error) {
	pageNum := index >> b.exp
	offset := index & b.r

	bufferSlice, err := b.get(pageNum)
	if err != nil {
		return nil, erroring.NewInvalidArgumentError("index", erroring.IndexOutOfRange, err)
	}
	if index < 0 || index >= len(bufferSlice) {
		return nil, erroring.NewInvalidArgumentError("index", erroring.IndexOutOfRange, nil)
	}
	v := bufferSlice[offset]
	return v, nil
}

func (b *FastObjectBuffer) ModifyEntry(index int, value interface{}) error {
	pageNum := index >> b.exp
	bufferSlice, err := b.get(pageNum)
	if err != nil {
		return erroring.NewInvalidArgumentError("index", erroring.IndexOutOfRange, err)
	}

	offset := index & b.r
	if index < 0 || index >= len(bufferSlice) {
		return erroring.NewInvalidArgumentError("index", erroring.IndexOutOfRange, nil)
	}
	bufferSlice[offset] = value
	return b.set(pageNum, bufferSlice)
}

func (b *FastObjectBuffer) GetSize() int {
	return b.size
}

func (b *FastObjectBuffer) SetSize(size int) {
	b.size = size
}

func (b *FastObjectBuffer) Append(value interface{}) error {
	if b.size < b.capacity {
		pageNum := b.size >> b.exp
		bufferSlice, err := b.get(pageNum)
		if err != nil {
			return err
		}
		bufferSlice = append(bufferSlice, value)
		return b.set(pageNum, bufferSlice)
	} else {
		b.size++
		b.capacity += b.pageSize
		b.buffer = append(b.buffer, []interface{}{value})
		return nil
	}
}

func (b *FastObjectBuffer) Clear() {
	b.size = 0
}

func (b *FastObjectBuffer) get(index int) ([]interface{}, error) {
	if index < 0 || index >= b.size {
		return nil, erroring.NewInvalidArgumentError("index", erroring.IndexOutOfRange, nil)
	}
	return b.buffer[index], nil
}

func (b *FastObjectBuffer) set(index int, value []interface{}) error {
	if index < 0 || index >= b.size {
		return erroring.NewInvalidArgumentError("index", erroring.IndexOutOfRange, nil)
	}
	b.buffer[index] = value
	return nil
}
