package reader

import (
	"unicode"

	"github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"
)

type Utf8Reader struct {
	xmlDoc    []byte
	offset    int
	endOffset int
}

func NewUtf8Reader(xmlDoc []byte, offset, endOffset int) (Utf8Reader, error) {
	if xmlDoc == nil {
		return Utf8Reader{}, erroring.NewInvalidArgumentError("xmlDoc", erroring.CannotBeNil, nil)
	}
	if offset < 0 {
		return Utf8Reader{}, erroring.NewInvalidArgumentError("offset", erroring.IndexOutOfRange, nil)
	}
	if endOffset < 0 || endOffset > len(xmlDoc) {
		return Utf8Reader{}, erroring.NewInvalidArgumentError("endOffset", erroring.IndexOutOfRange, nil)
	}
	return Utf8Reader{
		xmlDoc:    xmlDoc,
		offset:    offset,
		endOffset: endOffset,
	}, nil
}

func (r *Utf8Reader) GetChar() (uint32, error) {
	if r.offset >= r.endOffset {
		return 0, erroring.NewEOFError(erroring.XmlIncomplete)
	}
	ch := r.xmlDoc[r.offset]
	r.offset++
	if !r.isUTF8(ch) {
		return 0, erroring.NewParseError("invalid ASCII character", "", nil)
	}
	return uint32(ch), nil
}

func (r *Utf8Reader) GetLongChar(offset int32) (uint64, error) {
	ch := r.xmlDoc[offset]
	if ch == byte('\r') && r.xmlDoc[offset+1] == byte('\n') {
		return (2 << 32) | '\n', nil
	}
	return (1 << 32) | uint64(ch), nil
}

func (r *Utf8Reader) SkipChar(ch uint32) bool {
	if ch == uint32(r.xmlDoc[r.offset]) {
		r.offset++
		return true
	} else {
		return false
	}
}

func (r *Utf8Reader) SkipCharSeq(seq string) bool {
	for _, ch := range seq {
		if !r.SkipChar(uint32(ch)) {
			return false
		}
	}
	return true
}

func (r *Utf8Reader) Decode(offset int32) (uint32, error) {
	return uint32(r.xmlDoc[offset]), nil
}

func (r *Utf8Reader) GetOffset() int {
	return r.offset
}

func (r *Utf8Reader) SetOffset(offset int) {
	r.offset = offset
}

// isUTF8 function to validate if a character is a valid ASCII character
func (r *Utf8Reader) isUTF8(ch uint8) bool {
	return uint32(ch) <= unicode.MaxRune
}
