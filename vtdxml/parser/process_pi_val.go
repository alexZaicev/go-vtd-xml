package parser

import "github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"

func (p *VtdParser) processPiVal() (State, error) {
	if !p.xmlChar.IsSpaceChar(p.currentChar) {
		return StateInvalid, erroring.NewParseError("invalid termination sequence", p.fmtLine(), nil)
	}
	p.lastOffset = p.offset
	for {
		if err := p.nextChar(); err != nil {
			return StateInvalid, err
		}
		if !p.xmlChar.IsValidChar(p.currentChar) {
			return StateInvalid, erroring.NewParseError("invalid char in PI value", p.fmtLine(), nil)
		}
		if p.currentChar == '?' && p.skipChar('>') {
			break
		}
	}
	p.length1 = p.offset - p.lastOffset - (p.increment << 1)
	if err := p.writeVtdWithLengthCheck(TokenPiVal, "PI value too long (>0xFFFF"); err != nil {
		return StateInvalid, err
	}
	p.lastOffset = p.offset
	if err := p.nextCharAfterWs(); err != nil {
		return StateInvalid, err
	}
	return p.getNextProcessStateFromChar(p.currentChar)
}
