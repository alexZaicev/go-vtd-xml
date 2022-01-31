package parser

import (
	"github.com/alexZaicev/go-vtd-xml/vtdxml/buffer"
	"github.com/alexZaicev/go-vtd-xml/vtdxml/reader"
)

func (p *VtdParser) init() error {
	bufInt, err := buffer.NewFastIntBuffer([]buffer.FastIntBufferOption{
		buffer.WithFastIntBufferPageSize(DefaultNsBufferSize),
	}...)
	if err != nil {
		return err
	}
	p.nsBuffer1 = bufInt

	bufLong, err := buffer.NewFastLongBuffer([]buffer.FastLongBufferOption{
		buffer.WithFastLongBufferPageSize(DefaultNsBufferSize),
	}...)
	if err != nil {
		return err
	}
	p.nsBuffer2 = bufLong

	bufLong, err = buffer.NewFastLongBuffer([]buffer.FastLongBufferOption{
		buffer.WithFastLongBufferPageSize(DefaultNsBufferSize),
	}...)
	if err != nil {
		return err
	}
	p.nsBuffer3 = bufLong

	r, err := reader.NewUtf8Reader(p.xmlDoc, p.offset, p.endOffset)
	if err != nil {
		return err
	}
	p.reader = r

	if p.shallowDepth {
		if err := p.initWithShallowDepth(); err != nil {
			return err
		}
	} else {
		if err := p.initWithoutShallowDepth(); err != nil {
			return err
		}
	}

	return nil
}

func (p *VtdParser) initWithShallowDepth() error {
	vtdSize, l1size, l2size, l3size := 0, DefaultDepth3L1BufferSize, DefaultDepth3L2BufferSize,
		DefaultDepth3L3BufferSize
	if p.docLength <= 1024 {
		vtdSize = 6
		l1size = 5
		l2size = 5
		l3size = 5
	} else if p.docLength <= 4096 {
		vtdSize = 7
		l1size = 6
		l2size = 6
		l3size = 6
	} else if p.docLength <= 1024*16 {
		vtdSize = 8
		l1size = 7
		l2size = 7
		l3size = 7
	} else if p.docLength <= 1024*16*4 {
		vtdSize = 11
	} else if p.docLength <= 1024*256 {
		vtdSize = 12
	} else {
		vtdSize = 15
	}

	b, err := buffer.NewFastLongBuffer([]buffer.FastLongBufferOption{
		buffer.WithFastLongBufferPageSize(vtdSize),
	}...)
	if err != nil {
		return err
	}
	p.vtdBuffer = b

	b, err = buffer.NewFastLongBuffer([]buffer.FastLongBufferOption{
		buffer.WithFastLongBufferPageSize(l1size),
	}...)
	if err != nil {
		return err
	}
	p.l1Buffer = b

	b, err = buffer.NewFastLongBuffer([]buffer.FastLongBufferOption{
		buffer.WithFastLongBufferPageSize(l2size),
	}...)
	if err != nil {
		return err
	}
	p.l2Buffer = b

	b, err = buffer.NewFastLongBuffer([]buffer.FastLongBufferOption{
		buffer.WithFastLongBufferPageSize(l3size),
	}...)
	if err != nil {
		return err
	}
	p.l3Buffer = b

	return nil
}

func (p *VtdParser) initWithoutShallowDepth() error {
	vtdSize, l1size, l2size, l3size, l4size, l5size := 0, DefaultDepth5L1BufferSize, DefaultDepth5L2BufferSize,
		DefaultDepth5L3BufferSize, DefaultDepth5L4BufferSize, DefaultDepth5L5BufferSize
	if p.docLength <= 1024 {
		vtdSize = 6
		l1size = 5
		l2size = 5
		l3size = 5
		l4size = 5
		l5size = 5
	} else if p.docLength <= 4096 {
		vtdSize = 7
		l1size = 6
		l2size = 6
		l3size = 6
		l4size = 6
		l5size = 6
	} else if p.docLength <= 1024*16 {
		vtdSize = 8
		l1size = 7
		l2size = 7
		l3size = 7
		l4size = 7
		l5size = 7
	} else if p.docLength <= 1024*16*4 {
		vtdSize = 11
		l2size = 8
		l3size = 8
		l4size = 8
		l5size = 8
	} else if p.docLength <= 1024*256 {
		vtdSize = 12
		l1size = 8
		l2size = 9
		l3size = 9
		l4size = 9
		l5size = 9
	} else {
		vtdSize = 15
	}

	b, err := buffer.NewFastLongBuffer([]buffer.FastLongBufferOption{
		buffer.WithFastLongBufferPageSize(vtdSize),
	}...)
	if err != nil {
		return err
	}
	p.vtdBuffer = b

	b, err = buffer.NewFastLongBuffer([]buffer.FastLongBufferOption{
		buffer.WithFastLongBufferPageSize(l1size),
	}...)
	if err != nil {
		return err
	}
	p.l1Buffer = b

	b, err = buffer.NewFastLongBuffer([]buffer.FastLongBufferOption{
		buffer.WithFastLongBufferPageSize(l2size),
	}...)
	if err != nil {
		return err
	}
	p.l2Buffer = b

	b, err = buffer.NewFastLongBuffer([]buffer.FastLongBufferOption{
		buffer.WithFastLongBufferPageSize(l3size),
	}...)
	if err != nil {
		return err
	}
	p.l3Buffer = b

	b, err = buffer.NewFastLongBuffer([]buffer.FastLongBufferOption{
		buffer.WithFastLongBufferPageSize(l4size),
	}...)
	if err != nil {
		return err
	}
	p.l4Buffer = b

	b, err = buffer.NewFastLongBuffer([]buffer.FastLongBufferOption{
		buffer.WithFastLongBufferPageSize(l5size),
	}...)
	if err != nil {
		return err
	}
	p.l5Buffer = b

	return nil
}
