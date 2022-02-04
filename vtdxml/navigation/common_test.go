package navigation

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/alexZaicev/go-vtd-xml/vtdxml/buffer"
	"github.com/alexZaicev/go-vtd-xml/vtdxml/common"
	"github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"
	"github.com/stretchr/testify/assert"
)

var vtdBufferValues = []int64{
	-537919488, -1611661305, -1343225841, -1611661284, -1343225825, -1611661265, -1343225805, 44070659424314, 805316696,
	1073741915, 805316751, 1073741974, 805316815, 1073742039, 4547635927056624, 9051261324230921, 13554843771732268,
	5778118429290529084, 13554843771732342, 5778118429290529158, 13554835181797824, 5778118394930790862, 13554852361667070,
	5778118386340856336, 9051265619198558, 13510824651915906, 808452749, 1076888254, 18014407099417299, 22518015316722413,
	27021640713896717, 31525218866430775, 5796132767670469437, 18014407099417534, 22518015316722648, 27021640713896952,
	31525218866431010, 5796132767670469672, 18014437164188841, 5782621947313521843, 18014437164188888, 5782621938723587298,
	18014424279287045, 5782621973083325708, 18014419984319796, 5782622046097769786, 13510833241851269, 808453519, 1076888989,
	18014458639025607, 22518019611690477, 5787125546940892659, 22518028201625112, 5787125645725140512, 22518058266396254,
	5787125538350958189,
}

func Test_VtdNav_GetTokenType_Success(t *testing.T) {
	nav := getNav(t)

	testCases := []struct {
		name           string
		tokenID        int
		expectedResult int32
	}{
		{
			name:           "Token type at index 0",
			tokenID:        0,
			expectedResult: 15,
		},
		{
			name:           "Token type at index 8",
			tokenID:        8,
			expectedResult: 0,
		},
		{
			name:           "Token type at index 39",
			tokenID:        39,
			expectedResult: 5,
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			token, err := nav.GetTokenType(tc.tokenID)
			assert.Equal(t, tc.expectedResult, token)
			assert.Nil(t, err)
		})
	}

}

func Test_VtdNav_GetTokenType_Error(t *testing.T) {
	nav := getNav(t)

	testCases := []struct {
		name           string
		tokenID        int
		expectedErrMsg string
	}{
		{
			name:           "Token index out of range (9999)",
			tokenID:        9999,
			expectedErrMsg: "invalid argument index: array index out of range",
		},
		{
			name:           "Token index out of range (-1)",
			tokenID:        -1,
			expectedErrMsg: "invalid argument index: array index out of range",
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			token, err := nav.GetTokenType(tc.tokenID)
			assert.Equal(t, int32(0), token)
			if assert.EqualError(t, err, tc.expectedErrMsg) {
				assert.IsType(t, erroring.InvalidArgumentErrorType, err)
			}
			assert.NoError(t, errors.Unwrap(err))
		})
	}
}

func Test_VtdNav_GetCurrentIndex_Success(t *testing.T) {
	nav := getNav(t)

	idx, err := nav.GetCurrentIndex()
	assert.Nil(t, err)
	assert.Equal(t, int32(0), idx)

	nav.context[0] = 0
	idx, err = nav.GetCurrentIndex()
	assert.Nil(t, err)
	assert.Equal(t, int32(7), idx)

	nav.atTerminal = true
	idx, err = nav.GetCurrentIndex()
	assert.Nil(t, err)
	assert.Equal(t, int32(0), idx)

	nav.atTerminal = false
	nav.context[0] = 2
	idx, err = nav.GetCurrentIndex()
	assert.Nil(t, err)
	assert.Equal(t, int32(-1), idx)
}

func Test_VtdNav_GetCurrentIndex_Error(t *testing.T) {
	nav := getNav(t)

	testCases := []struct {
		name           string
		tokenID        int32
		expectedErrMsg string
	}{
		{
			name:           "Context[0] value out of range (9999)",
			tokenID:        9999,
			expectedErrMsg: "an internal error occurred: array index out of range",
		},
		{
			name:           "Context[0] value out of range (-5)",
			tokenID:        -5,
			expectedErrMsg: "an internal error occurred: array index out of range",
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			nav.context[0] = tc.tokenID

			idx, err := nav.GetCurrentIndex()
			assert.Equal(t, int32(0), idx)
			if assert.EqualError(t, err, tc.expectedErrMsg) {
				assert.IsType(t, erroring.InternalErrorType, err)
			}
			assert.NoError(t, errors.Unwrap(err))
		})
	}
}

func Test_VtdNav_GetTokenOffset_Success(t *testing.T) {
	nav := getNav(t)

	testCases := []struct {
		name           string
		tokenID        int
		expectedResult int32
	}{
		{
			name:           "Token index at index 0",
			tokenID:        0,
			expectedResult: 535822336,
		},
		{
			name:           "Token index at index 8",
			tokenID:        8,
			expectedResult: 805316696,
		},
		{
			name:           "Token index at index 39",
			tokenID:        39,
			expectedResult: 1203,
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			token, err := nav.GetTokenOffset(tc.tokenID)
			assert.Equal(t, tc.expectedResult, token)
			assert.Nil(t, err)
		})
	}
}

func Test_VtdNav_GetTokenOffset_Error(t *testing.T) {
	nav := getNav(t)

	testCases := []struct {
		name           string
		tokenID        int
		expectedErrMsg string
	}{
		{
			name:           "Token index out of range (9999)",
			tokenID:        9999,
			expectedErrMsg: "invalid argument index: array index out of range",
		},
		{
			name:           "Token index out of range (-1)",
			tokenID:        -1,
			expectedErrMsg: "invalid argument index: array index out of range",
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			token, err := nav.GetTokenOffset(tc.tokenID)
			assert.Equal(t, int32(0), token)
			if assert.EqualError(t, err, tc.expectedErrMsg) {
				assert.IsType(t, erroring.InvalidArgumentErrorType, err)
			}
			assert.NoError(t, errors.Unwrap(err))
		})
	}
}

func Test_VtdNav_GetTokenDepth_Success(t *testing.T) {
	nav := getNav(t)

	testCases := []struct {
		name           string
		tokenID        int
		expectedResult int32
	}{
		{
			name:           "Token depth at index 0",
			tokenID:        0,
			expectedResult: -1,
		},
		{
			name:           "Token depth at index 8",
			tokenID:        8,
			expectedResult: 0,
		},
		{
			name:           "Token depth at index 39",
			tokenID:        39,
			expectedResult: 4,
		},
		{
			name:           "Token depth at index 55",
			tokenID:        55,
			expectedResult: 5,
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			token, err := nav.GetTokenDepth(tc.tokenID)
			assert.Equal(t, tc.expectedResult, token)
			assert.Nil(t, err)
		})
	}
}

func Test_VtdNav_GetTokenDepth_Error(t *testing.T) {
	nav := getNav(t)

	testCases := []struct {
		name           string
		tokenID        int
		expectedErrMsg string
	}{
		{
			name:           "Token depth out of range (9999)",
			tokenID:        9999,
			expectedErrMsg: "invalid argument index: array index out of range",
		},
		{
			name:           "Token depth out of range (-1)",
			tokenID:        -1,
			expectedErrMsg: "invalid argument index: array index out of range",
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			token, err := nav.GetTokenDepth(tc.tokenID)
			assert.Equal(t, int32(0), token)
			if assert.EqualError(t, err, tc.expectedErrMsg) {
				assert.IsType(t, erroring.InvalidArgumentErrorType, err)
			}
			assert.NoError(t, errors.Unwrap(err))
		})
	}
}

func Test_VtdNav_GetTokenLength_Success(t *testing.T) {
	nav := getNav(t)

	testCases := []struct {
		name           string
		tokenID        int
		expectedResult int32
	}{
		{
			name:           "Token length at index 0",
			tokenID:        0,
			expectedResult: 1048575,
		},
		{
			name:           "Token length at index 6",
			tokenID:        6,
			expectedResult: 1048575,
		},
		{
			name:           "Token length at index 7",
			tokenID:        7,
			expectedResult: 327701,
		},
		{
			name:           "Token length at index 11",
			tokenID:        11,
			expectedResult: 0,
		},
		{
			name:           "Token length at index 23",
			tokenID:        23,
			expectedResult: 15,
		},
		{
			name:           "Token length at index 45",
			tokenID:        45,
			expectedResult: 29,
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			token, err := nav.GetTokenLength(tc.tokenID)
			assert.Equal(t, tc.expectedResult, token)
			assert.Nil(t, err)
		})
	}
}

func Test_VtdNav_GetTokenLength_Error(t *testing.T) {
	nav := getNav(t)

	testCases := []struct {
		name           string
		tokenID        int
		expectedErrMsg string
	}{
		{
			name:           "Token length out of range (9999)",
			tokenID:        9999,
			expectedErrMsg: "invalid argument index: array index out of range",
		},
		{
			name:           "Token length out of range (-1)",
			tokenID:        -1,
			expectedErrMsg: "invalid argument index: array index out of range",
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			token, err := nav.GetTokenLength(tc.tokenID)
			assert.Equal(t, int32(0), token)
			if assert.EqualError(t, err, tc.expectedErrMsg) {
				assert.IsType(t, erroring.InvalidArgumentErrorType, err)
			}
			assert.NoError(t, errors.Unwrap(err))
		})
	}
}

func getNav(t *testing.T) *VtdNav {
	p := filepath.Join("..", "..", "testdata", "xml_valid", "mirs.095.001.01.golden")
	path, err := filepath.Abs(p)
	assert.Nil(t, err)

	g, err := ioutil.ReadFile(path)
	assert.Nil(t, err)
	assert.NotNil(t, g)

	vtdBuffer, err := buffer.NewFastLongBuffer([]buffer.FastLongBufferOption{
		buffer.WithFastLongBufferPageSize(7),
	}...)
	assert.Nil(t, err)
	assert.NotNil(t, vtdBuffer)
	for _, i := range vtdBufferValues {
		assert.Nil(t, vtdBuffer.Append(i))
	}

	l1Buffer, err := buffer.NewFastLongBuffer([]buffer.FastLongBufferOption{
		buffer.WithFastLongBufferPageSize(6),
	}...)
	assert.Nil(t, err)
	assert.NotNil(t, l1Buffer)

	l2Buffer, err := buffer.NewFastLongBuffer([]buffer.FastLongBufferOption{
		buffer.WithFastLongBufferPageSize(6),
	}...)
	assert.Nil(t, err)
	assert.NotNil(t, l2Buffer)

	l3Buffer, err := buffer.NewFastLongBuffer([]buffer.FastLongBufferOption{
		buffer.WithFastLongBufferPageSize(6),
	}...)
	assert.Nil(t, err)
	assert.NotNil(t, l3Buffer)

	nav, err := NewVtdNav(7, 0, int32(len(g)), 7,
		common.FormatUtf8, true, g,
		vtdBuffer, l1Buffer, l2Buffer, l3Buffer)

	assert.Nil(t, err)
	assert.NotNil(t, nav)
	return nav
}
