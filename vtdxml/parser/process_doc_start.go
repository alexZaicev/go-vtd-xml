package parser

import "github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"

func (p *VtdParser) processDocStart() (State, error) {
	ch, err := p.reader.GetChar()
	if err != nil {
		return StateInvalid, err
	}
	if p.currentChar == '<' {
		p.lastOffset = p.offset
		if p.reader.SkipChar('?') &&
			(p.reader.SkipChar('x') || p.reader.SkipChar('X')) &&
			(p.reader.SkipChar('m') || p.reader.SkipChar('M')) &&
			(p.reader.SkipChar('l') || p.reader.SkipChar('L')) {
			if p.reader.SkipChar('?') ||
				p.reader.SkipChar('\t') ||
				p.reader.SkipChar('\n') ||
				p.reader.SkipChar('\r') {
				if err := p.nextCharAfterWs(); err != nil {
					return StateInvalid, err
				}
				p.lastOffset = p.offset
				return StateDecAttrName, nil
			} else if p.reader.SkipChar('?') {
				return StateInvalid, erroring.NewParseError("premature ending", p.fmtLine(), nil)
			}
		}
		p.offset = p.lastOffset
		return StateLtSeen, nil
	} else if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
		if err := p.nextCharAfterWs(); err != nil {
			return StateInvalid, nil
		}
		if p.currentChar == '<' {
			return StateLtSeen, nil
		}
	}
	return StateInvalid, nil
}
