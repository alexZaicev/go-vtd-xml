package buffer

import (
	"github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"
)

const (
	DefaultLongPageSize = 1024
)

type FastLongBufferOption func(*FastLongBuffer)

type LongBuffer interface {
	Buffer
	LongAt(index int) (int64, error)
	ModifyEntry(index int, value int64) error
	Lower32At(index int) (int32, error)
	Upper32At(index int) (int32, error)
	ToLongArray() ([]int64, error)
	Append(value int64) error
}

type FastLongBuffer struct {
	buffer   [][]int64
	capacity int
	size     int
	exp      int
	r        int
	pageSize int
}

func WithFastLongBufferPageSize(size int) FastLongBufferOption {
	return func(b *FastLongBuffer) {
		b.pageSize = 1 << size
		b.exp = size
		b.r = b.pageSize - 1
	}
}

func NewFastLongBuffer(opts ...FastLongBufferOption) (*FastLongBuffer, error) {
	b := &FastLongBuffer{
		capacity: 0,
		size:     0,
		pageSize: DefaultLongPageSize,
		exp:      10,
		r:        DefaultIntPageSize - 1,
		buffer:   make([][]int64, 0, DefaultLongPageSize),
	}

	for _, opt := range opts {
		opt(b)
	}
	if b.pageSize == 0 || b.r < 0 {
		return nil, erroring.NewInvalidArgumentError("size", erroring.InvalidBufferPageSize, nil)
	}

	return b, nil
}

func (b *FastLongBuffer) LongAt(index int) (int64, error) {
	pageNum := index >> b.exp
	offset := index & b.r

	bufferSlice, err := b.get(pageNum)
	if err != nil {
		return 0, erroring.NewInvalidArgumentError("index", erroring.IndexOutOfRange, err)
	}
	if index < 0 || index >= len(bufferSlice) {
		return 0, erroring.NewInvalidArgumentError("index", erroring.IndexOutOfRange, nil)
	}
	v := bufferSlice[offset]
	return v, nil
}

func (b *FastLongBuffer) ModifyEntry(index int, value int64) error {
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

func (b *FastLongBuffer) GetSize() int {
	return b.size
}

func (b *FastLongBuffer) SetSize(size int) {
	b.size = size
}

// Lower32At function return lower 32 bit of the int64 at the index
func (b *FastLongBuffer) Lower32At(index int) (int32, error) {
	valueLong, err := b.LongAt(index)
	if err != nil {
		return 0, erroring.NewInvalidArgumentError("index", erroring.IndexOutOfRange, err)
	}
	valueInt := int32(valueLong)
	return valueInt, nil
}

// Upper32At function return upper 32 bit of the int64 at the index
func (b *FastLongBuffer) Upper32At(index int) (int32, error) {
	valueLong, err := b.LongAt(index)
	if err != nil {
		return 0, erroring.NewInvalidArgumentError("index", erroring.IndexOutOfRange, err)
	}
	valueLong = valueLong & (0xFFFFFFF << 32)
	valueLong = valueLong >> 32
	valueInt := int32(valueLong)
	return valueInt, nil
}

// ToLongArray function converts 2D buffer into int64 slice
func (b *FastLongBuffer) ToLongArray() ([]int64, error) {
	// if buffer is empty return empty slice
	if b.size < 1 {
		return []int64{}, nil
	}
	size := b.size
	var intArray []int64

	for i := 0; size > 0; i++ {
		// get buffer page slice
		buffer, err := b.get(i)
		if err != nil {
			// if error occurs stop iteration and return error
			return nil, err
		}
		// load-in buffer page into int64 slice
		for j := range buffer {
			v := buffer[j]
			intArray = append(intArray, v)
		}
		// subtract buffer size with read page size
		size -= b.pageSize
	}
	return intArray, nil
}

func (b *FastLongBuffer) Append(value int64) error {
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
		b.buffer = append(b.buffer, []int64{value})
		return nil
	}
}

func (b *FastLongBuffer) Clear() {
	b.size = 0
}

func (b *FastLongBuffer) get(index int) ([]int64, error) {
	if index < 0 || index >= b.size {
		return nil, erroring.NewInvalidArgumentError("index", erroring.IndexOutOfRange, nil)
	}
	return b.buffer[index], nil
}

func (b *FastLongBuffer) set(index int, value []int64) error {
	if index < 0 || index >= b.size {
		return erroring.NewInvalidArgumentError("index", erroring.IndexOutOfRange, nil)
	}
	b.buffer[index] = value
	return nil
}
