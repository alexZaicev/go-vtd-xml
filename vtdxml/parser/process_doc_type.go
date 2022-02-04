package parser

import (
	"github.com/alexZaicev/go-vtd-xml/vtdxml/common"
	"github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"
)

func (p *VtdParser) processDocType() (State, error) {
	z := 1
	for {
		ch, err := p.getChar()
		if err != nil {
			return StateInvalid, err
		}
		if !p.xmlChar.IsValidChar(ch) {
			return StateInvalid, erroring.NewParseError("invalid char in DOCTYPE", p.fmtLine(), nil)
		}
		if ch == '>' {
			z--
		} else if ch == '<' {
			z++
		}

		if z == 0 {
			break
		}
	}

	p.length1 = p.offset - p.lastOffset - p.increment
	if err := p.writeVtdWithLengthCheck(common.TokenDtdVal, "DTD value too long >0xFFFFF"); err != nil {
		return StateInvalid, err
	}
	if err := p.nextCharAfterWs(); err != nil {
		return StateInvalid, err
	}
	if p.currentChar != '<' {
		return StateInvalid, erroring.NewParseError(erroring.InvalidChar, p.fmtLine(), nil)
	}
	return StateLtSeen, nil
}
