package navigation

import (
	"errors"
	"testing"

	"github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"
	"github.com/stretchr/testify/assert"
)

func Test_VtdNav_ToStringAtRange_Success(t *testing.T) {
	nav := getNav(t)

	testCases := []struct {
		name           string
		offset         int32
		length         int32
		expectedResult string
	}{
		{
			name:           "Range [0, 50]",
			offset:         0,
			length:         55,
			expectedResult: "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?>",
		},
		{
			name:           "Range [240, 242]",
			offset:         240,
			length:         2,
			expectedResult: "Sw",
		},
		{
			name:           "Range [328, 338]",
			offset:         328,
			length:         10,
			expectedResult: "ius33,o=sw",
		},
		{
			name:           "Range [0, 300]",
			offset:         0,
			length:         300,
			expectedResult: "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?>\n\n<SwInt:ExchangeRequest xmlns:Sw=\"urn:swift:snl:ns.Sw\"\n\n                       xmlns:SwInt=\"urn:swift:snl:ns.SwInt\"\n\n                       xmlns:SwSec=\"urn:swift:snl:ns.SwSec\">\n\n    <SwInt:Request>\n\n        <SwInt:RequestHeader>\n\n            <",
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			token, err := nav.ToStringAtRange(tc.offset, tc.length)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResult, token)
		})
	}
}

func Test_VtdNav_ToStringAtRange_Error(t *testing.T) {
	nav := getNav(t)

	testCases := []struct {
		name           string
		offset         int32
		length         int32
		expectedErrMsg string
	}{
		{
			name:           "Negative index",
			offset:         -1,
			length:         55,
			expectedErrMsg: "invalid argument index: array index out of range",
		},
		{
			name:           "Length too big",
			offset:         0,
			length:         999999,
			expectedErrMsg: "invalid argument index: array index out of range",
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			str, err := nav.ToStringAtRange(tc.offset, tc.length)
			assert.Equal(t, "", str)
			if assert.EqualError(t, err, tc.expectedErrMsg) {
				assert.IsType(t, erroring.InvalidArgumentErrorType, err)
			}
			assert.NoError(t, errors.Unwrap(err))
		})
	}
}

func Test_VtdNav_ToRawStringAtRange_Success(t *testing.T) {
	nav := getNav(t)

	testCases := []struct {
		name           string
		offset         int32
		length         int32
		expectedResult string
	}{
		{
			name:           "Range [0, 50]",
			offset:         0,
			length:         55,
			expectedResult: "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?>",
		},
		{
			name:           "Range [240, 242]",
			offset:         240,
			length:         2,
			expectedResult: "Sw",
		},
		{
			name:           "Range [328, 338]",
			offset:         328,
			length:         10,
			expectedResult: "ius33,o=sw",
		},
		{
			name:           "Range [0, 300]",
			offset:         0,
			length:         300,
			expectedResult: "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?>\n\n<SwInt:ExchangeRequest xmlns:Sw=\"urn:swift:snl:ns.Sw\"\n\n                       xmlns:SwInt=\"urn:swift:snl:ns.SwInt\"\n\n                       xmlns:SwSec=\"urn:swift:snl:ns.SwSec\">\n\n    <SwInt:Request>\n\n        <SwInt:RequestHeader>\n\n            <",
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			token, err := nav.ToRawStringAtRange(tc.offset, tc.length)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResult, token)
		})
	}
}

func Test_VtdNav_ToRawStringAtRange_Error(t *testing.T) {
	nav := getNav(t)

	testCases := []struct {
		name           string
		offset         int32
		length         int32
		expectedErrMsg string
	}{
		{
			name:           "Negative index",
			offset:         -1,
			length:         55,
			expectedErrMsg: "invalid argument index: array index out of range",
		},
		{
			name:           "Length too big",
			offset:         0,
			length:         999999,
			expectedErrMsg: "invalid argument index: array index out of range",
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			str, err := nav.ToRawStringAtRange(tc.offset, tc.length)
			assert.Equal(t, "", str)
			if assert.EqualError(t, err, tc.expectedErrMsg) {
				assert.IsType(t, erroring.InvalidArgumentErrorType, err)
			}
			assert.NoError(t, errors.Unwrap(err))
		})
	}
}

func Test_VtdNav_ToRawStringAtIndex_Success(t *testing.T) {
	nav := getNav(t)

	testCases := []struct {
		name           string
		index          int
		expectedResult string
	}{
		{
			name:           "Range [0, 50]",
			index:          7,
			expectedResult: "SwInt:ExchangeRequest",
		},
		{
			name:           "Range [240, 242]",
			index:          14,
			expectedResult: "SwInt:Request",
		},
		{
			name:           "Range [328, 338]",
			index:          19,
			expectedResult: "ou=xxx,o=citius33,o=swift",
		},
		{
			name:           "Range [0, 300]",
			index:          25,
			expectedResult: "AppHdr",
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			token, err := nav.ToRawStringAtIndex(tc.index)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResult, token)
		})
	}
}

func Test_VtdNav_ToRawStringAtIndex_Error(t *testing.T) {
	nav := getNav(t)

	testCases := []struct {
		name           string
		index          int
		expectedErrMsg string
	}{
		{
			name:           "Not existing index",
			index:          5,
			expectedErrMsg: "invalid argument index: array index out of range",
		},
		{
			name:           "Negative index",
			index:          -1,
			expectedErrMsg: "invalid argument index: array index out of range",
		},
		{
			name:           "Index too big",
			index:          99999,
			expectedErrMsg: "invalid argument index: array index out of range",
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			str, err := nav.ToRawStringAtIndex(tc.index)
			assert.Equal(t, "", str)
			if assert.EqualError(t, err, tc.expectedErrMsg) {
				assert.IsType(t, erroring.InvalidArgumentErrorType, err)
			}
			assert.NoError(t, errors.Unwrap(err))
		})
	}
}
