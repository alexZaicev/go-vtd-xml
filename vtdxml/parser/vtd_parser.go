package parser

import (
	"github.com/alexZaicev/go-vtd-xml/vtdxml/buffer"
	"github.com/alexZaicev/go-vtd-xml/vtdxml/common"
	"github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"
	"github.com/alexZaicev/go-vtd-xml/vtdxml/reader"
)

const (
	DefaultNsBufferSize = 6
	DefaultDepth        = -1
	DefaultLcDepth      = 3
	DefaultIncrement    = 1

	DefaultBufferReuse = false

	DefaultDepth3L1BufferSize = 8
	DefaultDepth3L2BufferSize = 9
	DefaultDepth3L3BufferSize = 11

	DefaultDepth5L1BufferSize = 7
	DefaultDepth5L2BufferSize = 9
	DefaultDepth5L3BufferSize = 11
	DefaultDepth5L4BufferSize = 11
	DefaultDepth5L5BufferSize = 11

	DefaultTagArraySize  = 256
	DefaultAttrArraySize = 256
	DefaultUrArraySize   = 256
)

type Option func(*VtdParser)

// VtdParser VTD generator implementation supporting build-in entities only.
// Handles DTD parsing, but does not resolve declared entities.
type VtdParser struct {
	xmlDoc                                                              []byte
	offset, docOffset, lastOffset, endOffset                            int
	length, length1, length2, docLength                                 int
	depth, vtdDepth, lcDepth, lastDepth                                 int
	increment                                                           int
	rootIndex, lastL1Index, lastL2Index, lastL3Index, lastL4Index       int
	attrCount, prefixedAttCount                                         int
	attrNameSlice, prefixUrlSlice                                       []int
	currentChar, lastChar                                               uint32
	currentElementRecord                                                int64
	nsAware, defaultNs, isNs                                            bool
	singleByteEncoding, bomDetected, mustUtf8, shallowDepth, helper, ws bool
	isXml                                                               bool
	bufferReuse                                                         bool
	encoding                                                            common.FormatEncoding
	xmlChar                                                             *common.XmlChar
	vtdBuffer, l1Buffer, l2Buffer, l3Buffer, l4Buffer, l5Buffer         buffer.LongBuffer
	nsBuffer1                                                           buffer.IntBuffer
	nsBuffer2, nsBuffer3                                                buffer.LongBuffer
	reader                                                              reader.Reader
	tagStack, prefixedAttrNameSlice                                     []int64
}

func WithXmlDoc(xmlDoc []byte) Option {
	return func(p *VtdParser) {
		p.xmlDoc = xmlDoc
		p.offset, p.docOffset = 0, 0
		p.length = len(p.xmlDoc)
		p.docLength = p.length
		p.endOffset = p.length
	}
}

func WithXmlDocCustomOffset(xmlDoc []byte, offset, length int) Option {
	return func(p *VtdParser) {
		p.xmlDoc = xmlDoc
		p.offset, p.docOffset = offset, offset
		p.length, p.docLength = length, length
		p.endOffset = offset + length
	}
}

func WithNameSpaceAware(nsAware bool) Option {
	return func(p *VtdParser) {
		p.nsAware = nsAware
	}
}

func WithBufferReuse(br bool) Option {
	return func(p *VtdParser) {
		p.bufferReuse = br
	}
}

func WithLcDepth(depth int) Option {
	return func(p *VtdParser) {
		if depth == 3 {
			p.shallowDepth = true
		} else if depth == 5 {
			p.shallowDepth = false
		}
	}
}

// NewVtdParser function creates a new instance of VTD parser with
// options provided, error otherwise
func NewVtdParser(opts ...Option) (*VtdParser, error) {
	g := &VtdParser{
		xmlChar:               common.NewXmlChar(),
		singleByteEncoding:    true,
		shallowDepth:          true,
		depth:                 DefaultDepth,
		lcDepth:               DefaultLcDepth,
		bufferReuse:           DefaultBufferReuse,
		increment:             DefaultIncrement,
		tagStack:              make([]int64, DefaultTagArraySize, DefaultTagArraySize),
		attrNameSlice:         make([]int, DefaultAttrArraySize, DefaultAttrArraySize),
		prefixedAttrNameSlice: make([]int64, DefaultAttrArraySize, DefaultAttrArraySize),
		prefixUrlSlice:        make([]int, DefaultAttrArraySize, DefaultAttrArraySize),
	}

	for _, opt := range opts {
		opt(g)
	}

	if g.xmlDoc == nil {
		return nil, erroring.NewInvalidArgumentError("xmlDoc", erroring.CannotBeNil, nil)
	}
	if len(g.xmlDoc) == 0 {
		return nil, erroring.NewInvalidArgumentError("xmlDoc", "document cannot be empty", nil)
	}
	if g.offset < 0 {
		return nil, erroring.NewInvalidArgumentError("offset", erroring.IndexOutOfRange, nil)
	}
	if g.length == 0 || g.offset+g.length > len(g.xmlDoc) {
		return nil, erroring.NewInvalidArgumentError("length", erroring.InvalidSliceLength, nil)
	}
	// perform additional initialization
	if err := g.init(); err != nil {
		return nil, err
	}

	return g, nil
}
