package reader

import (
	"fmt"
	"github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"
	"io"
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
	return AsciiReader{
		xmlDoc:    xmlDoc,
		offset:    offset,
		endOffset: endOffset,
	}, nil
}

func (r *AsciiReader) GetChar() (int32, error) {
	if r.offset >= r.endOffset {
		return 0, io.EOF
	}
	ch := r.xmlDoc[r.offset]
	r.offset++
	if ch < 0 {
		return 0, erroring.NewParseError(fmt.Sprintf("invalid ASCII chararacter: %d", ch), nil)
	}
	return int32(ch), nil
}

func (r *AsciiReader) GetLongChar(offset int32) (int64, error) {
	ch := r.xmlDoc[offset]
	if ch == byte('\r') && r.xmlDoc[offset+1] == byte('\n') {
		return (2 << 32) | '\n', nil
	}
	return (1 << 32) | int64(ch), nil
}

func (r *AsciiReader) SkipChar(ch int32) (bool, error) {
	if ch == int32(r.xmlDoc[r.offset]) {
		r.offset++
		return true, nil
	} else {
		return false, nil
	}
}

func (r *AsciiReader) Decode(offset int32) (string, error) {
	return fmt.Sprintf("%c", r.xmlDoc[offset]), nil
}
