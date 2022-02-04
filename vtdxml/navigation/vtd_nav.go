package navigation

import (
	"github.com/alexZaicev/go-vtd-xml/vtdxml/buffer"
	"github.com/alexZaicev/go-vtd-xml/vtdxml/common"
)

type Nav interface {
	GetCurrentIndex() (int32, error)
	GetTokenType(index int) (int32, error)
	GetTokenOffset(index int) (int32, error)
	GetTokenLength(index int) (int32, error)
	GetTokenDepth(index int) (int32, error)

	ToStringAtIndex(index int) (string, error)
	ToStringAtRange(offset, length int32) (string, error)

	ToRawStringAtIndex(index int) (string, error)
	ToRawStringAtRange(index, length int32) (string, error)

	ToElement(dir Direction) (bool, error)
}

type VtdNav struct {
	context                                 []int32
	rootIndex, offset, length, depth        int32
	encoding                                common.FormatEncoding
	nsAware                                 bool
	atTerminal                              bool
	ln                                      int32
	xmlChar                                 *common.XmlChar
	xmlBuffer                               buffer.ByteBuffer
	vtdBuffer, l1Buffer, l2Buffer, l3Buffer buffer.LongBuffer
	l1index, l2index, l3index               int
	l2lower, l2upper, l3lower, l3upper      int32
}

func NewVtdNav(
	rootIndex, offset, length, depth int32,
	encoding common.FormatEncoding,
	nsAware bool,
	bytes []byte,
	vtdBuffer, l1Buffer, l2Buffer, l3Buffer buffer.LongBuffer,
) (*VtdNav, error) {
	// TODO validate arguments
	n := &VtdNav{
		rootIndex: rootIndex,
		offset:    offset,
		length:    length,
		depth:     depth + 1,
		encoding:  encoding,
		nsAware:   nsAware,
		vtdBuffer: vtdBuffer,
		l1Buffer:  l1Buffer,
		l2Buffer:  l2Buffer,
		l3Buffer:  l3Buffer,
		context:   make([]int32, 0, depth),
		xmlChar:   common.NewXmlChar(),
	}

	for i := 0; i < int(depth); i++ {
		n.context = append(n.context, -1)
	}

	xmlBuffer, err := buffer.NewUniByteBuffer(bytes)
	if err != nil {
		return nil, err
	}
	n.xmlBuffer = xmlBuffer

	return n, nil
}
