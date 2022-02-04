package parser

import "github.com/alexZaicev/go-vtd-xml/vtdxml/common"

func (p *VtdParser) processElementTail() (State, error) {
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
			if p.skipChar('/') {
				if p.helper {
					p.length1 = p.offset - p.lastOffset - (p.increment << 1)
					if p.singleByteEncoding {
						if err := p.writeVtdText(common.TokenCharacterData, p.lastOffset, p.length1, p.depth); err != nil {
							return StateInvalid, err
						}
					} else {
						if err := p.writeVtdText(common.TokenCharacterData, p.lastOffset>>1, p.length1>>1,
							p.depth); err != nil {
							return StateInvalid, err
						}
					}
				}
				return StateTagEnd, nil
			}
			return StateLtSeen, nil
		} else if p.xmlChar.IsContentChar(p.currentChar) {
			return StateText, nil
		} else {
			if err := p.handleOtherTextChar(p.currentChar); err != nil {
				return StateInvalid, err
			}
			return StateText, nil
		}
	}
	return StateDocEnd, nil
}
