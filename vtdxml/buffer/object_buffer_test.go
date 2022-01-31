package buffer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	defaultObjectSize       = 32
	defaultObjectBufferSize = 1024
)

func Test_NewFastObjectBuffer_Success(t *testing.T) {
	buffer, err := NewFastObjectBuffer()
	assert.NotNil(t, buffer)
	assert.Nil(t, err)
	assert.Equal(t, 0, buffer.size)
	assert.Equal(t, 0, buffer.capacity)
	assert.Equal(t, 10, buffer.exp)
	assert.Equal(t, defaultObjectBufferSize, buffer.pageSize)
	assert.Equal(t, defaultObjectBufferSize-1, buffer.r)
	assert.NotNil(t, buffer.buffer)

	expectedPageSize := 1 << defaultObjectSize

	buffer, err = NewFastObjectBuffer(WithFastObjectBufferPageSize(defaultObjectSize))
	assert.Nil(t, err)
	assert.NotNil(t, buffer)
	assert.Equal(t, 0, buffer.size)
	assert.Equal(t, 0, buffer.capacity)
	assert.Equal(t, defaultObjectSize, buffer.exp)
	assert.Equal(t, expectedPageSize, buffer.pageSize)
	assert.Equal(t, expectedPageSize-1, buffer.r)
	assert.NotNil(t, buffer.buffer)
}

func Test_NewFastObjectBuffer_InvalidArgument(t *testing.T) {
	size := 2048

	buffer, err := NewFastIntBuffer(WithFastIntBufferPageSize(size))
	assert.EqualError(t, err, "invalid argument size: invalid buffer page size")
	assert.Nil(t, buffer)
}

func Test_FastObjectBuffer_ObjectAt_Success(t *testing.T) {
	buffer := getInitializedFastObjectBuffer(t)
	for i := 0; i < defaultLongBufferSize; i++ {
		value, err := buffer.ObjectAt(i)
		assert.Nil(t, err)
		assert.Equal(t, int64(i), value)
	}

	for i := defaultLongBufferSize; i < defaultLongBufferSize*2; i++ {
		value, err := buffer.ObjectAt(i)
		assert.Nil(t, err)
		assert.Equal(t, fmt.Sprintf("%d", i), value)
	}

	for i := defaultLongBufferSize * 2; i < defaultLongBufferSize*3; i++ {
		value, err := buffer.ObjectAt(i)
		assert.Nil(t, err)
		assert.Equal(t, int32(i), value)
	}

	for i := defaultLongBufferSize * 3; i < defaultLongBufferSize*4; i++ {
		value, err := buffer.ObjectAt(i)
		assert.Nil(t, err)
		assert.Equal(t, byte(i), value)
	}
}

func Test_FastObjectBuffer_ObjectAt_InvalidArgument(t *testing.T) {
	buffer := getInitializedFastObjectBuffer(t)

	testCases := []struct {
		name           string
		index          int
		expectedErrMsg string
	}{
		{
			name:           "negative index",
			index:          -1,
			expectedErrMsg: "invalid argument index: array index out of range",
		},
		{
			name:           "index of array size",
			index:          defaultLongBufferSize * 4,
			expectedErrMsg: "invalid argument index: array index out of range",
		},
		{
			name:           "index greater than array size",
			index:          defaultLongBufferSize*4 + 1,
			expectedErrMsg: "invalid argument index: array index out of range",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			b, err := buffer.ObjectAt(tc.index)
			assert.EqualError(t, err, tc.expectedErrMsg)
			assert.Nil(t, b)
		})
	}
}

func Test_FastObjectBuffer_ModifyEntry_Success(t *testing.T) {
	buffer := getInitializedFastObjectBuffer(t)

	assert.Nil(t, buffer.ModifyEntry(10, int64(100)))
	v, err := buffer.ObjectAt(10)
	assert.Nil(t, err)
	assert.Equal(t, int64(100), v)

	assert.Nil(t, buffer.ModifyEntry(10, int32(985)))
	v, err = buffer.ObjectAt(10)
	assert.Nil(t, err)
	assert.Equal(t, int32(985), v)

	assert.Nil(t, buffer.ModifyEntry(10, byte(8)))
	v, err = buffer.ObjectAt(10)
	assert.Nil(t, err)
	assert.Equal(t, byte(8), v)

	assert.Nil(t, buffer.ModifyEntry(10, "hello world"))
	v, err = buffer.ObjectAt(10)
	assert.Nil(t, err)
	assert.Equal(t, "hello world", v)
}

func Test_FastObjectBuffer_ModifyEntry_InvalidArgument(t *testing.T) {
	buffer := getInitializedFastObjectBuffer(t)

	testCases := []struct {
		name           string
		index          int
		value          int64
		expectedErrMsg string
	}{
		{
			name:           "negative index",
			index:          -1,
			value:          int64(1234),
			expectedErrMsg: "invalid argument index: array index out of range",
		},
		{
			name:           "index of array size",
			index:          defaultLongBufferSize * 4,
			value:          int64(1234),
			expectedErrMsg: "invalid argument index: array index out of range",
		},
		{
			name:           "index greater than array size",
			index:          defaultLongBufferSize*4 + 1,
			value:          int64(1234),
			expectedErrMsg: "invalid argument index: array index out of range",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := buffer.ModifyEntry(tc.index, tc.value)
			assert.EqualError(t, err, tc.expectedErrMsg)
		})
	}
}

func Test_FastObjectBuffer_Clear_Success(t *testing.T) {
	buffer := getInitializedFastObjectBuffer(t)
	buffer.Clear()
	assert.Equal(t, 0, buffer.GetSize())
}

func getInitializedFastObjectBuffer(t *testing.T) *FastObjectBuffer {
	buffer, err := NewFastObjectBuffer(WithFastObjectBufferPageSize(defaultLongSize))
	assert.Nil(t, err)
	assert.NotNil(t, buffer)

	for i := 0; i < defaultObjectBufferSize; i++ {
		err := buffer.Append(int64(i))
		assert.Nil(t, err)
	}

	for i := defaultObjectBufferSize; i < defaultObjectBufferSize*2; i++ {
		err := buffer.Append(fmt.Sprintf("%d", i))
		assert.Nil(t, err)
	}

	for i := defaultObjectBufferSize * 2; i < defaultObjectBufferSize*3; i++ {
		err := buffer.Append(int32(i))
		assert.Nil(t, err)
	}

	for i := defaultObjectBufferSize * 3; i < defaultObjectBufferSize*4; i++ {
		err := buffer.Append(byte(i))
		assert.Nil(t, err)
	}

	// we expect our buffer to have exactly 1 page
	assert.Equal(t, 1, buffer.GetSize())

	return buffer
}
