package reader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const XML = "<?xml version=\"1.0\" encoding=\"uTf-8\" standalone=\"no\"?>\n<pre:Vehicle xmlns:pre='urn:example-org:Transport' type='car'>\n\t<seats> 4 </seats>\n\t<colour> White </colour>\n\t<engine>\n\t\t<petrol/>\n\t\t<capacity units='cc'>1598</capacity>\n\t</engine>\n</pre:Vehicle>"

func Test_Utf8Reader_NewUtf8Reader_Success(t *testing.T) {
	r, err := NewUtf8Reader([]byte(XML), 0, len(XML))
	assert.Nil(t, err)
	assert.NotNil(t, r)
}

func Test_Utf8Reader_NewUtf8Reader_InvalidArg(t *testing.T) {
	testCases := []struct {
		name           string
		docBytes       []byte
		offset         int
		endOffset      int
		expectedErrMsg string
	}{
		{
			name:           "nil doc bytes",
			expectedErrMsg: "invalid argument xmlDoc: cannot be nil",
		},
		{
			name:           "invalid offset",
			docBytes:       []byte(XML),
			offset:         -1,
			expectedErrMsg: "invalid argument offset: array index out of range",
		},
		{
			name:           "invalid end offset negative",
			docBytes:       []byte(XML),
			offset:         0,
			endOffset:      -1,
			expectedErrMsg: "invalid argument endOffset: array index out of range",
		},
		{
			name:           "invalid end offset bigger than doc length",
			docBytes:       []byte(XML),
			offset:         0,
			endOffset:      999999999,
			expectedErrMsg: "invalid argument endOffset: array index out of range",
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r, err := NewUtf8Reader(tc.docBytes, tc.offset, tc.endOffset)
			assert.Nil(t, r)
			assert.NotNil(t, err)
			assert.EqualError(t, err, tc.expectedErrMsg)
		})
	}
}

func Test_Utf8Reader_GetChar_Success(t *testing.T) {
	r, err := NewUtf8Reader([]byte(XML), 0, len(XML))
	assert.Nil(t, err)
	assert.NotNil(t, r)

	for i := 0; i < len(XML); i++ {
		ch, err := r.GetChar()
		assert.Nil(t, err)
		assert.True(t, ch > 0)
	}
}

func Test_Utf8Reader_GetChar_Failed(t *testing.T) {
	testCases := []struct {
		name           string
		docBytes       []byte
		offset         int
		endOffset      int
		expectedErrMsg string
	}{
		{
			name:           "premature EOF",
			expectedErrMsg: "premature EOF reached: XML document incomplete",
			docBytes:       []byte(XML),
			offset:         len(XML),
			endOffset:      len(XML),
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r, err := NewUtf8Reader(tc.docBytes, tc.offset, tc.endOffset)
			assert.NotNil(t, r)
			assert.Nil(t, err)

			ch, err := r.GetChar()
			assert.NotNil(t, err)
			assert.EqualError(t, err, tc.expectedErrMsg)
			assert.Equal(t, uint32(0), ch)
		})
	}
}

func Test_Utf8Reader_SkipChar_Success(t *testing.T) {
	r, err := NewUtf8Reader([]byte(XML), 0, len(XML))
	assert.Nil(t, err)
	assert.NotNil(t, r)

	assert.True(t, r.SkipChar('<'))
	assert.True(t, r.SkipChar('?'))
	assert.True(t, r.SkipChar('x'))
}

func Test_Utf8Reader_SkipCharSeq_Success(t *testing.T) {
	r, err := NewUtf8Reader([]byte(XML), 0, len(XML))
	assert.Nil(t, err)
	assert.NotNil(t, r)

	assert.True(t, r.SkipCharSeq("<?xml "))
}
