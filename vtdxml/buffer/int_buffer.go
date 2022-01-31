package buffer

import (
	"github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"
)

const (
	DefaultIntPageSize = 1024
)

type FastIntBufferOption func(*FastIntBuffer)

type IntBuffer interface {
	Buffer
	IntAt(index int) (int32, error)
	ModifyEntry(index int, value int32) error
	ToIntArray() ([]int32, error)
	Append(value int32) error
}

type FastIntBuffer struct {
	buffer   [][]int32
	capacity int
	size     int
	exp      int
	r        int
	pageSize int
}

func WithFastIntBufferPageSize(size int) FastIntBufferOption {
	return func(b *FastIntBuffer) {
		if size >= 0 {
			b.pageSize = 1 << size
			b.exp = size
			b.r = b.pageSize - 1
		}
	}
}

func NewFastIntBuffer(opts ...FastIntBufferOption) (*FastIntBuffer, error) {
	b := &FastIntBuffer{
		capacity: 0,
		size:     0,
		pageSize: DefaultIntPageSize,
		exp:      10,
		r:        DefaultIntPageSize - 1,
		buffer:   make([][]int32, 0, DefaultLongPageSize),
	}

	for _, opt := range opts {
		opt(b)
	}

	if b.pageSize == 0 || b.r < 0 {
		return nil, erroring.NewInvalidArgumentError("size", erroring.InvalidBufferPageSize, nil)
	}

	return b, nil
}

func (b *FastIntBuffer) IntAt(index int) (int32, error) {
	pageNum := index >> b.exp
	offset := index & b.r

	bufferSlice, err := b.get(pageNum)
	if err != nil {
		return 0, err
	}
	if index < 0 || index >= len(bufferSlice) {
		return 0, erroring.NewInvalidArgumentError("index", erroring.IndexOutOfRange, nil)
	}
	v := bufferSlice[offset]
	return v, nil
}

func (b *FastIntBuffer) ModifyEntry(index int, value int32) error {
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

func (b *FastIntBuffer) GetSize() int {
	return b.size
}

func (b *FastIntBuffer) SetSize(size int) {
	b.size = size
}

// ToIntArray function converts 2D buffer into int32 slice
func (b *FastIntBuffer) ToIntArray() ([]int32, error) {
	// if buffer is empty return empty slice
	if b.size < 1 {
		return []int32{}, nil
	}
	size := b.size
	var intArray []int32

	for i := 0; size > 0; i++ {
		// get buffer page slice
		buffer, err := b.get(i)
		if err != nil {
			// if error occurs stop iteration and return error
			return nil, err
		}
		// load-in buffer page into int32 slice
		for j := range buffer {
			v := buffer[j]
			intArray = append(intArray, v)
		}
		// subtract buffer size with read page size
		size -= b.pageSize
	}
	return intArray, nil
}

func (b *FastIntBuffer) Append(value int32) error {
	if b.size < b.capacity {
		pageNum := b.size >> b.exp
		bufferSlice, err := b.get(pageNum)
		if err != nil {
			return err
		}
		bufferSlice = append(bufferSlice, value)
		b.size++
		return b.set(pageNum, bufferSlice)
	} else {
		b.size++
		b.capacity += b.pageSize
		b.buffer = append(b.buffer, []int32{value})
		return nil
	}
}

func (b *FastIntBuffer) Clear() {
	b.size = 0
}

func (b *FastIntBuffer) get(index int) ([]int32, error) {
	if index < 0 || index >= b.size {
		return nil, erroring.NewInvalidArgumentError("index", erroring.IndexOutOfRange, nil)
	}
	return b.buffer[index], nil
}

func (b *FastIntBuffer) set(index int, value []int32) error {
	if index < 0 || index >= b.size {
		return erroring.NewInvalidArgumentError("index", erroring.IndexOutOfRange, nil)
	}
	b.buffer[index] = value
	return nil
}
