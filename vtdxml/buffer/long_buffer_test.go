package buffer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	defaultLongSize       = 32
	defaultLongBufferSize = 1024
)

func Test_NewFastLongBuffer_Success(t *testing.T) {
	buffer, err := NewFastLongBuffer()
	assert.NotNil(t, buffer)
	assert.Nil(t, err)
	assert.Equal(t, 0, buffer.size)
	assert.Equal(t, 0, buffer.capacity)
	assert.Equal(t, 10, buffer.exp)
	assert.Equal(t, defaultLongBufferSize, buffer.pageSize)
	assert.Equal(t, defaultLongBufferSize-1, buffer.r)
	assert.NotNil(t, buffer.buffer)
	assert.Equal(t, 0, buffer.buffer.Size())

	expectedPageSize := 1 << defaultLongSize

	buffer, err = NewFastLongBuffer(WithFastLongBufferPageSize(defaultLongSize))
	assert.Nil(t, err)
	assert.NotNil(t, buffer)
	assert.Equal(t, 0, buffer.size)
	assert.Equal(t, 0, buffer.capacity)
	assert.Equal(t, defaultLongSize, buffer.exp)
	assert.Equal(t, expectedPageSize, buffer.pageSize)
	assert.Equal(t, expectedPageSize-1, buffer.r)
	assert.NotNil(t, buffer.buffer)
	assert.Equal(t, 0, buffer.buffer.Size())
}

func Test_NewFastLongBuffer_InvalidArgument(t *testing.T) {
	size := 2048

	buffer, err := NewFastIntBuffer(WithFastIntBufferPageSize(size))
	assert.EqualError(t, err, "invalid argument size: invalid buffer page size")
	assert.NotNil(t, buffer)
}

func Test_FastLongBuffer_LongAt_Success(t *testing.T) {
	buffer := getInitializedFastLongBuffer(t)
	for i := 0; i < defaultLongBufferSize; i++ {
		value, err := buffer.LongAt(i)
		assert.Nil(t, err)
		assert.Equal(t, int64(i), value)
	}
}

func Test_FastLongBuffer_LongAt_InvalidArgument(t *testing.T) {
	buffer := getInitializedFastLongBuffer(t)

	testCases := []struct {
		name             string
		index            int
		expectedErrMsg   string
		expectedIntValue int64
	}{
		{
			name:             "negative index",
			index:            -1,
			expectedErrMsg:   "invalid argument index: array index out of range",
			expectedIntValue: 0,
		},
		{
			name:             "index of array size",
			index:            defaultLongBufferSize,
			expectedErrMsg:   "invalid argument index: array index out of range",
			expectedIntValue: 0,
		},
		{
			name:             "index greater than array size",
			index:            defaultLongBufferSize + 1,
			expectedErrMsg:   "invalid argument index: array index out of range",
			expectedIntValue: 0,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			b, err := buffer.LongAt(tc.index)
			assert.EqualError(t, err, tc.expectedErrMsg)
			assert.Equal(t, tc.expectedIntValue, b)
		})
	}
}

func Test_FastLongBuffer_ModifyEntry_Success(t *testing.T) {
	buffer := getInitializedFastLongBuffer(t)

	assert.Nil(t, buffer.ModifyEntry(10, int64(100)))
	v, err := buffer.LongAt(10)
	assert.Nil(t, err)
	assert.Equal(t, int64(100), v)

	assert.Nil(t, buffer.ModifyEntry(1000, int64(985)))
	v, err = buffer.LongAt(1000)
	assert.Nil(t, err)
	assert.Equal(t, int64(985), v)
}

func Test_FastLongBuffer_ModifyEntry_InvalidArgument(t *testing.T) {
	buffer := getInitializedFastLongBuffer(t)

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
			index:          defaultLongBufferSize,
			value:          int64(1234),
			expectedErrMsg: "invalid argument index: array index out of range",
		},
		{
			name:           "index greater than array size",
			index:          defaultLongBufferSize + 1,
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

func Test_FastLongBuffer_Clear_Success(t *testing.T) {
	buffer := getInitializedFastLongBuffer(t)
	buffer.Clear()
	assert.Equal(t, 0, buffer.GetSize())
}

func Test_FastLongBuffer_Lower32At_Success(t *testing.T) {
	buffer := getInitializedFastLongBuffer(t)

	testCases := []struct {
		name          string
		index         int
		value         int64
		expectedValue int32
	}{
		{
			name:          "value 100",
			index:         100,
			value:         int64(100),
			expectedValue: int32(100),
		},
		{
			name:          "value 1234567891234",
			index:         50,
			value:         int64(1234567891234),
			expectedValue: int32(1912277282),
		},
		{
			name:          "value 9685749585475",
			index:         950,
			value:         int64(9685749585475),
			expectedValue: int32(598332995),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := buffer.ModifyEntry(tc.index, tc.value)
			assert.Nil(t, err)

			value32, err := buffer.Lower32At(tc.index)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedValue, value32)
		})
	}
}

func Test_FastLongBuffer_Lower32At_InvalidArgument(t *testing.T) {
	buffer := getInitializedFastLongBuffer(t)

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
			index:          defaultLongBufferSize,
			value:          int64(1234),
			expectedErrMsg: "invalid argument index: array index out of range",
		},
		{
			name:           "index greater than array size",
			index:          defaultLongBufferSize + 1,
			value:          int64(1234),
			expectedErrMsg: "invalid argument index: array index out of range",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			value32, err := buffer.Lower32At(tc.index)
			assert.EqualError(t, err, tc.expectedErrMsg)
			assert.Equal(t, int32(0), value32)
		})
	}
}

func Test_FastLongBuffer_Upper32At_Success(t *testing.T) {
	buffer := getInitializedFastLongBuffer(t)

	testCases := []struct {
		name          string
		index         int
		value         int64
		expectedValue int32
	}{
		{
			name:          "value 100",
			index:         100,
			value:         int64(100),
			expectedValue: int32(0),
		},
		{
			name:          "value 1234567891234",
			index:         50,
			value:         int64(1234567891234),
			expectedValue: int32(287),
		},
		{
			name:          "value 9685749585475",
			index:         950,
			value:         int64(9685749585475),
			expectedValue: int32(2255),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := buffer.ModifyEntry(tc.index, tc.value)
			assert.Nil(t, err)

			value32, err := buffer.Upper32At(tc.index)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedValue, value32)
		})
	}
}

func Test_FastLongBuffer_Upper32At_InvalidArgument(t *testing.T) {
	buffer := getInitializedFastLongBuffer(t)

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
			index:          defaultLongBufferSize,
			value:          int64(1234),
			expectedErrMsg: "invalid argument index: array index out of range",
		},
		{
			name:           "index greater than array size",
			index:          defaultLongBufferSize + 1,
			value:          int64(1234),
			expectedErrMsg: "invalid argument index: array index out of range",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			value32, err := buffer.Upper32At(tc.index)
			assert.EqualError(t, err, tc.expectedErrMsg)
			assert.Equal(t, int32(0), value32)
		})
	}
}

func getInitializedFastLongBuffer(t *testing.T) FastLongBuffer {
	buffer, err := NewFastLongBuffer(WithFastLongBufferPageSize(defaultLongSize))
	assert.Nil(t, err)
	assert.NotNil(t, buffer)

	for i := 0; i < defaultLongBufferSize; i++ {
		err := buffer.Append(int64(i))
		assert.Nil(t, err)
	}

	// we expect our buffer to have exactly 1 page
	assert.Equal(t, 1, buffer.GetSize())

	longSlice, err := buffer.ToLongArray()
	assert.Nil(t, err)
	assert.Equal(t, defaultLongBufferSize, len(longSlice))
	for i := 0; i < defaultLongBufferSize; i++ {
		assert.Equal(t, int64(i), longSlice[i])
	}

	return buffer
}
