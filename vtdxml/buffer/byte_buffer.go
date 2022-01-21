package buffer

import "github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"

type ByteBuffer interface {
	ByteAt(index int) (byte, error)
	GetByteSlice(offset int, length int) ([]byte, error)
	GetBytes() []byte
	Size() int
}

type UniByteBuffer struct {
	buffer []byte
}

func NewUniByteBuffer(buffer []byte) (UniByteBuffer, error) {
	if buffer == nil {
		return UniByteBuffer{}, erroring.NewInvalidArgumentError("buffer", erroring.CannotBeNil, nil)
	}
	return UniByteBuffer{
		buffer: buffer,
	}, nil
}

func (b *UniByteBuffer) ByteAt(index int) (byte, error) {
	if index >= len(b.buffer) || index < 0 {
		return 0, erroring.NewInvalidArgumentError("index", erroring.IndexOutOfRange, nil)
	}
	return b.buffer[index], nil
}

func (b *UniByteBuffer) GetByteSlice(offset int, length int) ([]byte, error) {
	if offset < 0 || offset >= len(b.buffer) {
		return nil, erroring.NewInvalidArgumentError("offset", erroring.IndexOutOfRange, nil)
	}
	if length < 1 || offset+length >= len(b.buffer) {
		return nil, erroring.NewInvalidArgumentError("length", erroring.InvalidSliceLength, nil)
	}
	return b.buffer[offset:length], nil
}

func (b *UniByteBuffer) GetBytes() []byte {
	return b.buffer
}

func (b *UniByteBuffer) Size() int {
	return len(b.buffer)
}
