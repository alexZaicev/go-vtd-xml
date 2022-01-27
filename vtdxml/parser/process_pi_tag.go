package parser

import "github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"

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
	if err := p.writeVtdWithLengthCheck(TokenPiName, "PI name too long >0xFFFFF"); err != nil {
		return StateInvalid, err
	}
	if p.currentChar == '?' {
		if p.singleByteEncoding {
			if err := p.writeVtd(TokenPiVal, p.lastOffset, 0, p.depth); err != nil {
				return StateInvalid, err
			}
		} else {
			if err := p.writeVtd(TokenPiVal, p.lastOffset>>1, 0, p.depth); err != nil {
				return StateInvalid, err
			}
		}
		if p.reader.SkipChar('>') {
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
			}
			if p.xmlChar.IsContentChar(p.currentChar) {
				return StateText, nil
			}
			if p.currentChar == '&' {
				ch, err := p.entityIdentifier()
				if err != nil {
					return StateInvalid, err
				}
				if !p.xmlChar.IsValidChar(ch) {
					return StateInvalid, erroring.NewParseError(erroring.InvalidChar, p.fmtLine(), nil)
				}
				return StateText, nil
			}
			if p.currentChar == ']' {
				// skip all ] chars
				for p.reader.SkipChar(']') {
				}
				if p.reader.SkipChar('>') {
					return StateInvalid, erroring.NewParseError("]]> sequence in text content", p.fmtLine(), nil)
				}
				return StateText, nil
			}
			return StateInvalid, erroring.NewParseError(erroring.InvalidChar, p.fmtLine(), nil)
		}
	} else {
		return StateInvalid, erroring.NewParseError("invalid PI termination sequence", p.fmtLine(), nil)
	}
	return StatePiVal, nil
}
