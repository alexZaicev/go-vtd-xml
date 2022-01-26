package parser

import "github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"

func (p *VtdParser) processEndTag() (State, error) {
	sOffset := int(p.tagStack[p.depth])
	sLength := int(p.tagStack[p.depth] >> 32)

	p.lastOffset = p.offset
	p.offset = p.lastOffset + sLength
	if p.offset >= p.endOffset {
		return StateInvalid, erroring.NewEOFError(erroring.XmlIncomplete)
	}

	for i := 0; i < sLength; i++ {
		if p.xmlDoc[sOffset+i] != p.xmlDoc[p.lastOffset+i] {
			return StateInvalid, erroring.NewParseError("start/end tag mismatch", p.fmtLine(), nil)
		}
	}
	p.depth--
	if err := p.nextCharAfterWs(); err != nil {
		return StateInvalid, err
	}
	if p.currentChar != '>' {
		return StateInvalid, erroring.NewParseError("invalid char in ending", p.fmtLine(), nil)
	}

	if p.depth != -1 {
		p.lastOffset = p.offset
		if err := p.nextCharAfterWs(); err != nil {
			return StateInvalid, err
		}
		if p.currentChar == '<' {
			if p.ws {
				if err := p.recordWhiteSpace(); err != nil {
					return StateInvalid, err
				}
			}
			return StateLtSeen, nil
		} else if p.xmlChar.IsContentChar(p.currentChar) {
			return StateText, nil
		} else {
			if err := p.handleOtherTextChar(p.currentChar); err != nil {
				return StateInvalid, err
			}
			return StateLtSeen, nil
		}
	}
	return StateDocEnd, nil
}
