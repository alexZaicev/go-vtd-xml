package parser

import (
	"github.com/alexZaicev/go-vtd-xml/vtdxml/buffer"
	"github.com/alexZaicev/go-vtd-xml/vtdxml/common"
	"github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"
	"github.com/alexZaicev/go-vtd-xml/vtdxml/reader"
)

type Option func(*VtdParser)

// VtdParser VTD generator implementation supporting build-in entities only.
// Handles DTD parsing, but does not resolve declared entities.
type VtdParser struct {
	xmlDoc                                                              []byte
	offset, docOffset, lastOffset, endOffset                            int
	length, length1, length2, docLength                                 int
	depth, vtdDepth, lastDepth                                          int
	increment                                                           int
	rootIndex, lastL1Index, lastL2Index, lastL3Index, lastL4Index       int
	attrCount, prefixedAttCount                                         int
	attrNameSlice, prefixedAttrNameSlice, prefixUrlSlice                []int
	currentChar, lastChar                                               uint32
	currentElementRecord                                                int64
	nsAware, defaultNs, isNs                                            bool
	singleByteEncoding, bomDetected, mustUtf8, shallowDepth, helper, ws bool
	isXml                                                               bool
	encoding                                                            FormatEncoding
	xmlChar                                                             *common.XmlChar
	vtdBuffer, l1Buffer, l2Buffer, l3Buffer, l4Buffer, l5Buffer         buffer.LongBuffer
	nsBuffer1                                                           buffer.IntBuffer
	nsBuffer2, nsBuffer3                                                buffer.LongBuffer
	reader                                                              reader.Reader
	tagStack                                                            []int64
}

func WithXmlDoc(xmlDoc []byte) Option {
	return func(g *VtdParser) {
		g.xmlDoc = xmlDoc
		g.offset = 0
		g.length = len(g.xmlDoc)
		g.docLength = len(g.xmlDoc)
	}
}

func WithXmlDocCustomOffset(xmlDoc []byte, offset, length int) Option {
	return func(g *VtdParser) {
		g.xmlDoc = xmlDoc
		g.offset = offset
		g.length = length
		g.docLength = length
	}
}

func WithNameSpaceAware(nsAware bool) Option {
	return func(g *VtdParser) {
		g.nsAware = nsAware
	}
}

func NewVtdGen(opts ...Option) (VtdParser, error) {
	g := VtdParser{
		xmlChar:            common.NewXmlChar(),
		singleByteEncoding: true,
	}

	for _, opt := range opts {
		opt(&g)
	}

	if g.xmlDoc == nil {
		return VtdParser{}, erroring.NewInvalidArgumentError("xmlDoc", erroring.CannotBeNil, nil)
	}
	if len(g.xmlDoc) == 0 {
		return VtdParser{}, erroring.NewInvalidArgumentError("xmlDoc", "document cannot be empty", nil)
	}
	if g.offset < 0 {
		return VtdParser{}, erroring.NewInvalidArgumentError("offset", erroring.IndexOutOfRange, nil)
	}
	if g.length == 0 || g.offset+g.length > len(g.xmlDoc) {
		return VtdParser{}, erroring.NewInvalidArgumentError("length", erroring.InvalidSliceLength, nil)
	}
	g.init()

	return g, nil
}

func (p *VtdParser) init() {
	// TODO
}
