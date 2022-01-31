package buffer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	defaultIntSize       = 32
	defaultIntBufferSize = 1024
)

func Test_NewFastIntBuffer_Success(t *testing.T) {
	buffer, err := NewFastIntBuffer()
	assert.NotNil(t, buffer)
	assert.Nil(t, err)
	assert.Equal(t, 0, buffer.size)
	assert.Equal(t, 0, buffer.capacity)
	assert.Equal(t, 10, buffer.exp)
	assert.Equal(t, defaultIntBufferSize, buffer.pageSize)
	assert.Equal(t, defaultIntBufferSize-1, buffer.r)
	assert.NotNil(t, buffer.buffer)
	assert.Equal(t, 0, buffer.buffer.Size())

	expectedPageSize := 1 << defaultIntSize

	buffer, err = NewFastIntBuffer(WithFastIntBufferPageSize(defaultIntSize))
	assert.Nil(t, err)
	assert.NotNil(t, buffer)
	assert.Equal(t, 0, buffer.size)
	assert.Equal(t, 0, buffer.capacity)
	assert.Equal(t, defaultIntSize, buffer.exp)
	assert.Equal(t, expectedPageSize, buffer.pageSize)
	assert.Equal(t, expectedPageSize-1, buffer.r)
	assert.NotNil(t, buffer.buffer)
	assert.Equal(t, 0, buffer.buffer.Size())
}

func Test_NewFastIntBuffer_InvalidArgument(t *testing.T) {
	size := 2048

	buffer, err := NewFastIntBuffer(WithFastIntBufferPageSize(size))
	assert.EqualError(t, err, "invalid argument size: invalid buffer page size")
	assert.Nil(t, buffer)
}

func Test_FastIntBuffer_IntAt_Success(t *testing.T) {
	buffer := getInitializedFastIntBuffer(t)
	for i := 0; i < defaultIntBufferSize; i++ {
		value, err := buffer.IntAt(i)
		assert.Nil(t, err)
		assert.Equal(t, int32(i), value)
	}
}

func Test_FastIntBuffer_IntAt_InvalidArgument(t *testing.T) {
	buffer := getInitializedFastIntBuffer(t)

	testCases := []struct {
		name             string
		index            int
		expectedErrMsg   string
		expectedIntValue int32
	}{
		{
			name:             "negative index",
			index:            -1,
			expectedErrMsg:   "invalid argument index: array index out of range",
			expectedIntValue: 0,
		},
		{
			name:             "index of array size",
			index:            defaultIntBufferSize,
			expectedErrMsg:   "invalid argument index: array index out of range",
			expectedIntValue: 0,
		},
		{
			name:             "index greater than array size",
			index:            defaultIntBufferSize + 1,
			expectedErrMsg:   "invalid argument index: array index out of range",
			expectedIntValue: 0,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			b, err := buffer.IntAt(tc.index)
			assert.EqualError(t, err, tc.expectedErrMsg)
			assert.Equal(t, tc.expectedIntValue, b)
		})
	}
}

func Test_FastIntBuffer_ModifyEntry_Success(t *testing.T) {
	buffer := getInitializedFastIntBuffer(t)

	assert.Nil(t, buffer.ModifyEntry(10, int32(100)))
	v, err := buffer.IntAt(10)
	assert.Nil(t, err)
	assert.Equal(t, int32(100), v)

	assert.Nil(t, buffer.ModifyEntry(1000, int32(985)))
	v, err = buffer.IntAt(1000)
	assert.Nil(t, err)
	assert.Equal(t, int32(985), v)
}

func Test_FastIntBuffer_ModifyEntry_InvalidArgument(t *testing.T) {
	buffer := getInitializedFastIntBuffer(t)

	testCases := []struct {
		name           string
		index          int
		value          int32
		expectedErrMsg string
	}{
		{
			name:           "negative index",
			index:          -1,
			value:          int32(1234),
			expectedErrMsg: "invalid argument index: array index out of range",
		},
		{
			name:           "index of array size",
			index:          defaultIntBufferSize,
			value:          int32(1234),
			expectedErrMsg: "invalid argument index: array index out of range",
		},
		{
			name:           "index greater than array size",
			index:          defaultIntBufferSize + 1,
			value:          int32(1234),
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

func Test_FastIntBuffer_Clear_Success(t *testing.T) {
	buffer := getInitializedFastIntBuffer(t)
	buffer.Clear()
	assert.Equal(t, 0, buffer.GetSize())
}

func getInitializedFastIntBuffer(t *testing.T) *FastIntBuffer {
	buffer, err := NewFastIntBuffer(WithFastIntBufferPageSize(defaultIntSize))
	assert.Nil(t, err)
	assert.NotNil(t, buffer)

	for i := 0; i < defaultIntBufferSize; i++ {
		err := buffer.Append(int32(i))
		assert.Nil(t, err)
	}

	assert.Equal(t, defaultLongBufferSize, buffer.GetSize())

	intSlice, err := buffer.ToIntArray()
	assert.Nil(t, err)
	assert.Equal(t, defaultIntBufferSize, len(intSlice))
	for i := 0; i < defaultIntBufferSize; i++ {
		assert.Equal(t, int32(i), intSlice[i])
	}

	return buffer
}
