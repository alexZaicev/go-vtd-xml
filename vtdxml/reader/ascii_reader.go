package reader

import (
	"unicode"

	"github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"
)

type AsciiReader struct {
	xmlDoc    []byte
	offset    int
	endOffset int
}

func NewAsciiReader(xmlDoc []byte, offset, endOffset int) (AsciiReader, error) {
	if xmlDoc == nil {
		return AsciiReader{}, erroring.NewInvalidArgumentError("xmlDoc", erroring.CannotBeNil, nil)
	}
	if offset < 0 {
		return AsciiReader{}, erroring.NewInvalidArgumentError("offset", erroring.IndexOutOfRange, nil)
	}
	if endOffset < 0 || endOffset >= len(xmlDoc) {
		return AsciiReader{}, erroring.NewInvalidArgumentError("endOffset", erroring.IndexOutOfRange, nil)
	}
	return AsciiReader{
		xmlDoc:    xmlDoc,
		offset:    offset,
		endOffset: endOffset,
	}, nil
}

func (r *AsciiReader) GetChar() (uint32, error) {
	if r.offset >= r.endOffset {
		return 0, erroring.NewEOFError(erroring.XmlIncomplete)
	}
	ch := r.xmlDoc[r.offset]
	r.offset++
	if !r.isASCII(ch) {
		return 0, erroring.NewParseError("invalid ASCII character", "", nil)
	}
	return uint32(ch), nil
}

func (r *AsciiReader) GetLongChar(offset int32) (uint64, error) {
	ch := r.xmlDoc[offset]
	if ch == byte('\r') && r.xmlDoc[offset+1] == byte('\n') {
		return (2 << 32) | '\n', nil
	}
	return (1 << 32) | uint64(ch), nil
}

func (r *AsciiReader) SkipChar(ch uint32) bool {
	if ch == uint32(r.xmlDoc[r.offset]) {
		r.offset++
		return true
	} else {
		return false
	}
}

func (r *AsciiReader) SkipCharSeq(seq string) bool {
	for _, ch := range seq {
		if !r.SkipChar(uint32(ch)) {
			return false
		}
	}
	return true
}

func (r *AsciiReader) Decode(offset int32) (uint32, error) {
	return uint32(r.xmlDoc[offset]), nil
}

// isASCII function to validate if a character is a valid ASCII character
func (r *AsciiReader) isASCII(ch uint8) bool {
	return ch <= unicode.MaxASCII
}
