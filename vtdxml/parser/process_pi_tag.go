package parser

import (
	"github.com/alexZaicev/go-vtd-xml/vtdxml/common"
	"github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"
)

func (p *VtdParser) processPiTag() (State, error) {
	for {
		if err := p.nextChar(); err != nil {
			return StateInvalid, err
		}
		if !p.xmlChar.IsNameChar(p.currentChar) {
			break
		}
	}
	p.length1 = p.offset - p.lastOffset - p.increment
	if err := p.writeVtdWithLengthCheck(common.TokenPiName, "PI name too long >0xFFFFF"); err != nil {
		return StateInvalid, err
	}
	if p.currentChar == '?' {
		if p.singleByteEncoding {
			if err := p.writeVtd(common.TokenPiVal, p.lastOffset, 0, p.depth); err != nil {
				return StateInvalid, err
			}
		} else {
			if err := p.writeVtd(common.TokenPiVal, p.lastOffset>>1, 0, p.depth); err != nil {
				return StateInvalid, err
			}
		}
		if p.skipChar('>') {
			p.lastOffset = p.offset
			if err := p.nextCharAfterWs(); err != nil {
				return StateInvalid, err
			}
			return p.getNextProcessStateFromChar(p.currentChar)
		} else {
			return StateInvalid, erroring.NewParseError("invalid PI termination sequence", p.fmtLine(), nil)
		}
	}
	return StatePiVal, nil
}
