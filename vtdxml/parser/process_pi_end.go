package parser

import (
	"github.com/alexZaicev/go-vtd-xml/vtdxml/common"
	"github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"
)

func (p *VtdParser) processPiEnd() (State, error) {
	if err := p.nextChar(); err != nil {
		return StateInvalid, err
	}
	if !p.xmlChar.IsNameStartChar(p.currentChar) {
		return StateInvalid, erroring.NewParseError("invalid char in PI target", p.fmtLine(), nil)
	}
	if (p.currentChar == 'x' || p.currentChar == 'X') &&
		(p.skipChar('m') || p.skipChar('M')) &&
		(p.skipChar('l') || p.skipChar('L')) {
		if err := p.nextChar(); err != nil {
			return StateInvalid, err
		}
		if p.xmlChar.IsSpaceChar(p.currentChar) || p.currentChar == '?' {
			return StateInvalid, erroring.NewParseError("[xX][mM][lL] not a valid PI target", p.fmtLine(), nil)
		}
	}
	for {
		if !p.xmlChar.IsNameChar(p.currentChar) {
			break
		}
		if err := p.nextChar(); err != nil {
			return StateInvalid, err
		}
	}

	p.length1 = p.offset - p.lastOffset - p.increment
	if err := p.writeVtdWithLengthCheck(common.TokenPiName, "PI name too long (>0xFFFF)"); err != nil {
		return StateInvalid, err
	}
	p.lastOffset = p.offset
	if p.xmlChar.IsSpaceChar(p.currentChar) {
		if err := p.nextCharAfterWs(); err != nil {
			return StateInvalid, err
		}
		for {
			if !p.xmlChar.IsValidChar(p.currentChar) {
				return StateInvalid, erroring.NewParseError("invalid char in PI value", p.fmtLine(), nil)
			}
			if p.currentChar == '?' && p.skipChar('>') {
				break
			}
			if err := p.nextChar(); err != nil {
				return StateInvalid, err
			}
		}
		p.length1 = p.offset - p.lastOffset - (p.increment << 1)
		if err := p.writeVtdWithLengthCheck(common.TokenPiVal, "PI value too long (>0xFFFF)"); err != nil {
			return StateInvalid, err
		}
	} else {
		if p.singleByteEncoding {
			if err := p.writeVtd(common.TokenPiVal, p.lastOffset, 0, p.depth); err != nil {
				return StateInvalid, err
			}
		} else {
			if err := p.writeVtd(common.TokenPiVal, p.lastOffset>>1, 0, p.depth); err != nil {
				return StateInvalid, err
			}
		}
		if p.currentChar != '?' || p.skipChar('>') {
			return StateInvalid, erroring.NewParseError("invalid termination sequence", p.fmtLine(), nil)
		}
	}
	return StateDocEnd, nil
}
