package parser

import "github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"

func (p *VtdParser) processStartTag() (State, error) {
	for {
		ch, err := p.reader.GetChar()
		if err != nil {
			return StateInvalid, err
		}
		if !p.xmlChar.IsNameChar(ch) {
			break
		}
		if ch == ':' {
			p.length2 = p.offset - p.lastOffset - p.increment
			if p.nsAware && p.checkXmlnsPrefix(p.lastOffset, p.length2, true) {
				return StateInvalid, erroring.NewParseError("XMLNS cannot be element prefix", p.fmtLine(), err)
			}
		}
	}

	p.length1 = p.offset - p.lastOffset - p.increment
	if p.depth > maxDepth {
		return StateInvalid, erroring.NewParseError(erroring.MaximumDepthExceeded, p.fmtLine(), nil)
	}

	x := (p.length1 << 32) + p.lastOffset
	p.tagStack[p.depth] = int64(x)

	if p.depth > p.vtdDepth {
		p.vtdDepth = p.depth
	}

	if p.singleByteEncoding {
		if p.length2 > maxPrefixLength || p.length1 > maxQnameLength {
			return StateInvalid, erroring.NewParseError(erroring.TagPrefixQnameTooLong, p.fmtLine(), nil)
		}
		if p.shallowDepth {
			if err := p.writeVtd(TokenStartingTag, p.lastOffset, (p.length2<<1)|p.length1, p.depth); err != nil {
				return StateInvalid, err
			}
		} else {
			if err := p.writeVtdL5(TokenStartingTag, p.lastOffset, (p.length2<<1)|p.length1, p.depth); err != nil {
				return StateInvalid, err
			}
		}
	} else {
		if p.length2 > (maxPrefixLength<<1) || p.length1 > (maxQnameLength<<1) {
			return StateInvalid, erroring.NewParseError(erroring.TagPrefixQnameTooLong, p.fmtLine(), nil)
		}
		if p.shallowDepth {
			if err := p.writeVtd(TokenStartingTag, p.lastOffset>>1, (p.length2<<10)|(p.length1>>1),
				p.depth); err != nil {
				return StateInvalid, err
			}
		} else {
			if err := p.writeVtdL5(TokenStartingTag, p.lastOffset>>1, (p.length2<<10)|(p.length1>>1),
				p.depth); err != nil {
				return StateInvalid, err
			}
		}
	}

	if p.nsAware {
		if p.length2 != 0 {
			p.length2 += p.increment
			p.currentElementRecord = int64(((p.length2<<16)|p.length1)<<32) | int64(p.lastOffset)
		} else {
			p.currentElementRecord = 0
		}

		if p.depth <= p.nsBuffer1.GetSize()-1 {
			p.nsBuffer1.SetSize(p.depth)
			t, err := p.nsBuffer1.IntAt(p.depth - 1)
			if err != nil {
				return StateInvalid, err
			}
			p.nsBuffer2.SetSize(int(t + 1))
			p.nsBuffer3.SetSize(int(t + 1))
		}
	}
	p.length2 = 0
	if p.xmlChar.IsSpaceChar(p.currentChar) {
		if err := p.nextCharAfterWs(); err != nil {
			return StateInvalid, err
		}
		if p.xmlChar.IsNameStartChar(p.currentChar) {
			offset, err := p.getPrevOffset()
			if err != nil {
				return StateInvalid, err
			}
			p.lastOffset = offset
			return StateAttrName, nil
		}
	}
	p.helper = true
	if p.currentChar == '/' {
		p.depth--
		p.helper = false
		ch, err := p.reader.GetChar()
		if err != nil {
			return StateInvalid, err
		}
		p.currentChar = ch
	}
	if p.currentChar == '>' {
		if p.nsAware {
			if err := p.nsBuffer1.Append(int32(p.nsBuffer3.GetSize() - 1)); err != nil {
				return StateInvalid, err
			}
			if p.currentElementRecord != 0 {
				if err := p.qualifyElement(); err != nil {
					return StateInvalid, err
				}
			}
		}
		return p.processElementTail()
	}
	return StateInvalid, erroring.NewParseError("invalid char in starting tag", p.fmtLine(), nil)
}
