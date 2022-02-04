package parser

import (
	"github.com/alexZaicev/go-vtd-xml/vtdxml/common"
	"github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"
)

func (p *VtdParser) processStartComment() (State, error) {
	for {
		if err := p.nextChar(); err != nil {
			return StateInvalid, err
		}
		if !p.xmlChar.IsValidChar(p.currentChar) {
			return StateInvalid, erroring.NewParseError(erroring.InvalidChar, p.fmtLine(), nil)
		}
		if p.currentChar == '-' && p.skipChar('-') {
			p.length1 = p.offset - p.lastOffset - (p.increment << 1)
			break
		}
	}
	if err := p.nextChar(); err != nil {
		return StateInvalid, err
	}
	if p.currentChar != '>' {
		return StateInvalid, erroring.NewParseError("invalid terminating sequence", p.fmtLine(), nil)
	}
	if p.singleByteEncoding {
		if err := p.writeVtdText(common.TokenComment, p.lastOffset, p.length1, p.depth); err != nil {
			return StateInvalid, err
		}
	} else {
		if err := p.writeVtdText(common.TokenComment, p.lastOffset>>1, p.length1>>1, p.depth); err != nil {
			return StateInvalid, err
		}
	}
	p.lastOffset = p.offset
	if err := p.nextCharAfterWs(); err != nil {
		return StateInvalid, err
	}
	return p.getNextProcessStateFromChar(p.currentChar)
}
