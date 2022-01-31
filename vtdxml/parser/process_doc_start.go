package parser

import "github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"

func (p *VtdParser) processDocStart() (State, error) {
	if err := p.nextChar(); err != nil {
		return StateInvalid, err
	}
	if p.currentChar == '<' {
		p.lastOffset = p.offset
		if p.skipChar('?') &&
			(p.skipChar('x') || p.skipChar('X')) &&
			(p.skipChar('m') || p.skipChar('M')) &&
			(p.skipChar('l') || p.skipChar('L')) {
			if p.skipChar(' ') ||
				p.skipChar('\t') ||
				p.skipChar('\n') ||
				p.skipChar('\r') {
				if err := p.nextCharAfterWs(); err != nil {
					return StateInvalid, err
				}
				p.lastOffset = p.offset
				return StateDecAttrName, nil
			} else if p.skipChar('?') {
				return StateInvalid, erroring.NewParseError("premature ending", p.fmtLine(), nil)
			}
		}
		p.offset = p.lastOffset
		return StateLtSeen, nil
	} else if p.currentChar == ' ' || p.currentChar == '\t' || p.currentChar == '\n' || p.currentChar == '\r' {
		if err := p.nextCharAfterWs(); err != nil {
			return StateInvalid, nil
		}
		if p.currentChar == '<' {
			return StateLtSeen, nil
		}
	}
	return StateInvalid, nil
}
