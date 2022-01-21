package buffer

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func Test_NewUniByteBuffer_Success(t *testing.T) {
	bytes := make([]byte, 1024)
	rand.Read(bytes)

	buffer, err := NewUniByteBuffer(bytes)
	assert.Nil(t, err)
	assert.NotNil(t, buffer)
	assert.Equal(t, 1024, buffer.Size())
}

func Test_NewUniByteBuffer_InvalidArguments(t *testing.T) {
	buffer, err := NewUniByteBuffer(nil)
	assert.EqualError(t, err, "invalid argument buffer: cannot be nil")
	assert.NotNil(t, buffer)
}

func Test_UniByteBuffer_ByteAt_Success(t *testing.T) {
	buffer := getInitializedBuffer(t)

	b, err := buffer.ByteAt(0)
	assert.Nil(t, err)
	assert.Equal(t, byte(15), b)

	b, err = buffer.ByteAt(1)
	assert.Nil(t, err)
	assert.Equal(t, byte(23), b)

	b, err = buffer.ByteAt(2)
	assert.Nil(t, err)
	assert.Equal(t, byte(10), b)
}

func Test_UniByteBuffer_ByteAt_InvalidArgument(t *testing.T) {
	buffer := getInitializedBuffer(t)

	testCases := []struct {
		name              string
		index             int
		expectedErrMsg    string
		expectedByteValue byte
	}{
		{
			name:              "negative index",
			index:             -1,
			expectedErrMsg:    "invalid argument index: array index out of range",
			expectedByteValue: 0,
		},
		{
			name:              "index of array size",
			index:             1024,
			expectedErrMsg:    "invalid argument index: array index out of range",
			expectedByteValue: 0,
		},
		{
			name:              "index greater than array size",
			index:             1025,
			expectedErrMsg:    "invalid argument index: array index out of range",
			expectedByteValue: 0,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			b, err := buffer.ByteAt(tc.index)
			assert.EqualError(t, err, tc.expectedErrMsg)
			assert.Equal(t, tc.expectedByteValue, b)
		})
	}
}

func Test_UniByteBuffer_GetByteSlice_Success(t *testing.T) {
	buffer := getInitializedBuffer(t)
	byteSlice, err := buffer.GetByteSlice(0, 3)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(byteSlice))
	assert.Equal(t, byte(15), byteSlice[0])
	assert.Equal(t, byte(23), byteSlice[1])
	assert.Equal(t, byte(10), byteSlice[2])
}

func Test_UniByteBuffer_GetByteSlice_InvalidArgument(t *testing.T) {
	buffer := getInitializedBuffer(t)

	testCases := []struct {
		name           string
		offset         int
		length         int
		expectedErrMsg string
	}{
		{
			name:           "negative index",
			offset:         -1,
			length:         3,
			expectedErrMsg: "invalid argument offset: array index out of range",
		},
		{
			name:           "index of array size",
			offset:         1024,
			length:         3,
			expectedErrMsg: "invalid argument offset: array index out of range",
		},
		{
			name:           "index greater than array size",
			offset:         1025,
			length:         3,
			expectedErrMsg: "invalid argument offset: array index out of range",
		},
		{
			name:           "negative length",
			offset:         0,
			length:         -1,
			expectedErrMsg: "invalid argument length: invalid slice length",
		},
		{
			name:           "length of array size",
			offset:         0,
			length:         1024,
			expectedErrMsg: "invalid argument length: invalid slice length",
		},
		{
			name:           "length greater than array size",
			offset:         0,
			length:         1024,
			expectedErrMsg: "invalid argument length: invalid slice length",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			byteSlice, err := buffer.GetByteSlice(tc.offset, tc.length)
			assert.EqualError(t, err, tc.expectedErrMsg)
			assert.Nil(t, byteSlice)
		})
	}
}

func getInitializedBuffer(t *testing.T) UniByteBuffer {
	bytes := make([]byte, 1024)
	bytes[0] = 15
	bytes[1] = 23
	bytes[2] = 10

	buffer, err := NewUniByteBuffer(bytes)
	assert.Nil(t, err)
	assert.NotNil(t, buffer)
	assert.Equal(t, 1024, buffer.Size())
	assert.Equal(t, bytes, buffer.GetBytes())

	return buffer
}
