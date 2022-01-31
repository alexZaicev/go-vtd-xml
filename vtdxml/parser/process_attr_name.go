package parser

import (
	"fmt"

	"github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"
)

func (p *VtdParser) processAttrName() (State, error) {
	if p.currentChar == 'x' {
		if p.skipCharSeq("mlns") {
			if err := p.nextChar(); err != nil {
				return StateInvalid, err
			}

			if p.currentChar == '=' || p.xmlChar.IsSpaceChar(p.currentChar) {
				p.defaultNs, p.isNs = true, true
			} else if p.currentChar == ':' {
				p.defaultNs, p.isNs = false, true
			}
		}
	}
	for {
		if !p.xmlChar.IsNameChar(p.currentChar) {
			break
		}
		if p.currentChar == ':' {
			p.length2 = p.offset - p.lastOffset - p.increment
		}
		if err := p.nextChar(); err != nil {
			return StateInvalid, err
		}
	}

	offset, err := p.getPrevOffset()
	if err != nil {
		return StateInvalid, err
	}
	p.length1 = offset - p.lastOffset

	if p.isNs && p.nsAware && !p.defaultNs {
		if (p.increment == 1 && (p.length1-p.length2 == 6)) ||
			(p.increment == 2 && (p.length1-p.length2 == 12)) {
			byteOffset := p.lastOffset + p.length2 + p.increment
			if p.checkXmlnsPrefix(byteOffset, -1, false) {
				return StateInvalid, erroring.NewParseError("XMLNS as namespace cannot be re-declared",
					p.fmtCustomLine(byteOffset), nil)
			}
		}
		if (p.increment == 1 && (p.length1-p.length2 == 4)) ||
			(p.increment == 2 && (p.length1-p.length2 == 8)) {
			byteOffset := p.lastOffset + p.length2 + p.increment
			p.isXml = p.checkXmlPrefix(byteOffset, -1, false)
		}
	}
	if err := p.checkAttrUniqueness(); err != nil {
		return StateInvalid, err
	}

	tokenType := TokenAttrName
	errMsg := erroring.AttrNamePrefixQnameTooLong
	if p.isNs {
		tokenType = TokenAttrNs
		errMsg = erroring.AttrNsPrefixQnameTooLong
		if p.nsAware && p.length2 != 0 && !p.isXml {
			val := int64((p.length2<<16|p.length1)<<32 | p.lastOffset)
			if err := p.nsBuffer3.Append(val); err != nil {
				return StateInvalid, err
			}
		}
	}

	if p.singleByteEncoding {
		if p.length2 > maxPrefixLength || p.length1 > maxQnameLength {
			return StateInvalid, erroring.NewParseError(errMsg, p.fmtLine(), nil)
		}
		if err := p.writeVtd(tokenType, p.lastOffset, (p.length2<<11)|p.length1, p.depth); err != nil {
			return StateInvalid, err
		}
	} else {
		if p.length2 > maxPrefixLength<<1 || p.length1 > maxQnameLength<<1 {
			return StateInvalid, erroring.NewParseError(errMsg, p.fmtLine(), nil)
		}
		if err := p.writeVtd(tokenType, p.lastOffset>>1, (p.length2<<10)|(p.length1>>1), p.depth); err != nil {
			return StateInvalid, err
		}
	}

	p.length2 = 0
	if p.xmlChar.IsSpaceChar(p.currentChar) {
		if err := p.nextCharAfterWs(); err != nil {
			return StateInvalid, err
		}
	}
	if p.currentChar != '=' {
		return StateInvalid, erroring.NewParseError(erroring.InvalidChar, p.fmtLine(), nil)
	}
	if err := p.nextCharAfterWs(); err != nil {
		return StateInvalid, err
	}
	if p.currentChar != '"' && p.currentChar != '\'' {
		return StateInvalid, erroring.NewParseError(fmt.Sprintf("%s should be ' or \" ", erroring.InvalidChar), p.fmtLine(), nil)
	}
	p.lastOffset = p.offset
	return StateAttrVal, nil
}

func (p *VtdParser) checkAttrUniqueness() error {
	var unique, uniqual bool
	unique = true
	for i := 0; i < p.attrCount; i++ {
		uniqual = false
		prevLen := int(int32(p.attrNameSlice[i]))
		if p.length1 == prevLen {
			prevOffset := int(int32(p.attrNameSlice[i] >> 32))
			for j := 0; j < prevLen; j++ {
				if p.xmlDoc[prevOffset+j] != p.xmlDoc[p.lastOffset+j] {
					uniqual = true
					break
				}
			}
		} else {
			uniqual = true
		}
		unique = unique && uniqual
	}
	// TODO possibly simplify to if !unique {}
	if !unique && p.attrCount != 0 {
		return erroring.NewParseError(erroring.AttrNotUnique, p.fmtLine(), nil)
	}
	p.attrNameSlice[p.attrCount] = (p.lastOffset << 32) | p.length1
	p.attrCount++
	if p.nsAware && !p.isNs && p.length2 != 0 {
		isXml := p.checkXmlPrefix(p.lastOffset, -1, false)
		if (p.increment == 1 && p.length2 == 3 && isXml) || (p.increment == 2 && p.length2 == 6 && isXml) {
			return nil
		}
		p.prefixedAttrNameSlice[p.prefixedAttCount] = int64((p.lastOffset << 32) | (p.length2 << 16) | p.length1)
		p.prefixedAttCount++
	}
	return nil
}
