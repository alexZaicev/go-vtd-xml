package parser

import "github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"

func (p *VtdParser) processDocEnd() (State, error) {
	if err := p.nextCharAfterWs(); err != nil {
		return StateInvalid, err
	}
	if p.currentChar == '<' {
		if p.skipChar('?') {
			p.lastOffset = p.offset
			return StatePiEnd, nil
		} else if p.skipCharSeq("!--") {
			p.lastOffset = p.offset
			return StateEndComment, nil
		}
	}
	return StateInvalid, erroring.NewParseError("XML not terminated properly", p.fmtLine(), nil)
}
