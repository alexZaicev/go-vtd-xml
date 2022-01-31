package parser

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var validXmlTestFiles = []string{
	"mirs.095.001.01.no_snl",
	"mirs.095.001.01",
	"camt.004.001.08",
	"camt.053.001.08",
	"xml_opt1",
	"xml_opt2",
	"xml_opt3",
	"xml_opt4",
	"xml_opt5",
}

func Test_VtdParser_Parse_WithoutNsAware_Success(t *testing.T) {
	for _, name := range validXmlTestFiles {
		name := name
		t.Run(name, func(t *testing.T) {
			parser, err := NewVtdParser([]Option{
				WithXmlDoc(readTestData(t, name, true)),
			}...)
			assert.Nil(t, err)
			assert.NotNil(t, *parser)

			assert.Nil(t, parser.Parse())
		})
	}
}

func Test_VtdParser_Parse_WithNsAware_Success(t *testing.T) {
	for _, name := range validXmlTestFiles {
		name := name
		t.Run(name, func(t *testing.T) {
			parser, err := NewVtdParser([]Option{
				WithXmlDoc(readTestData(t, name, true)),
				WithNameSpaceAware(true),
			}...)
			assert.Nil(t, err)
			assert.NotNil(t, *parser)

			assert.Nil(t, parser.Parse())
		})
	}
}

func Test_VtdParser_Parse_WithoutNsAware_InvalidDocument(t *testing.T) {
	testCase := []struct {
		name           string
		filename       string
		expectedErrMsg string
	}{
		{
			name:           "Illegal comments",
			filename:       "illegal_comments",
			expectedErrMsg: "a parse error occurred: invalid terminating sequence",
		},
		{
			name:           "Illegal declarative attribute name (version)",
			filename:       "illegal_dec_attr_version_name",
			expectedErrMsg: "a parse error occurred: declaration should be version",
		},
		{
			name:           "Illegal declarative attribute value (version)",
			filename:       "illegal_dec_attr_version_val",
			expectedErrMsg: "a parse error occurred: invalid version detected (supported 1.0 or 1.1)",
		},
		{
			name:           "Illegal declarative attribute name (encoding)",
			filename:       "illegal_dec_attr_encoding_name",
			expectedErrMsg: "a parse error occurred: declaration should be encoding",
		},
		{
			name:           "Illegal declarative attribute value (encoding)",
			filename:       "illegal_dec_attr_encoding_val",
			expectedErrMsg: "a parse error occurred: invalid document encoding",
		},
		{
			name:           "Illegal declarative attribute name (standalone)",
			filename:       "illegal_dec_attr_standalone_name",
			expectedErrMsg: "a parse error occurred: declaration should be standalone",
		},
		{
			name:     "Illegal declarative attribute value (standalone)",
			filename: "illegal_dec_attr_standalone_val",
			expectedErrMsg: "a parse error occurred: invalid value for attribute standalone (" +
				"valid options are yes or no)",
		},
	}
	for _, tc := range testCase {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			parser, err := NewVtdParser([]Option{
				WithXmlDoc(readTestData(t, tc.filename, false)),
				WithNameSpaceAware(false),
			}...)
			assert.Nil(t, err)
			assert.NotNil(t, *parser)

			err = parser.Parse()
			assert.NotNil(t, err)
			assert.EqualError(t, err, tc.expectedErrMsg)
		})
	}
}

func readTestData(t *testing.T, name string, valid bool) []byte {
	dir := "xml_invalid"
	if valid {
		dir = "xml_valid"
	}
	p := filepath.Join("..", "..", "testdata", dir, fmt.Sprintf("%s.golden", name))
	path, err := filepath.Abs(p)
	assert.Nil(t, err)

	g, err := ioutil.ReadFile(path)
	assert.Nil(t, err)
	assert.NotNil(t, g)
	return g
}
