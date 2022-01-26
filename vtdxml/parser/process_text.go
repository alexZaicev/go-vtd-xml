package parser

import "github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"

func (p *VtdParser) processText() (State, error) {
	if p.depth < 0 {
		return StateInvalid, erroring.NewParseError("text content at the wrong place", p.fmtLine(), nil)
	}
	for {
		if err := p.nextChar(); err != nil {
			return StateInvalid, err
		}
		if !p.xmlChar.IsContentChar(p.currentChar) {
			if p.currentChar == '<' {
				break
			}
			if err := p.handleOtherTextChar(p.currentChar); err != nil {
				return StateInvalid, err
			}
		}
	}
	p.length1 = p.offset - p.increment - p.lastOffset
	if p.singleByteEncoding {
		if err := p.writeVtdText(TokenCharacterData, p.lastOffset, p.length1, p.depth); err != nil {
			return StateInvalid, err
		}
	} else {
		if err := p.writeVtdText(TokenCharacterData, p.lastOffset>>1, p.length1>>1, p.depth); err != nil {
			return StateInvalid, err
		}
	}

	return StateLtSeen, nil
}
