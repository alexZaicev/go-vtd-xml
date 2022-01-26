package parser

import "github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"

func (p *VtdParser) processDocType() (State, error) {
	z := 1
	for {
		ch, err := p.reader.GetChar()
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

	if p.singleByteEncoding {
		if p.length1 > maxTokenLength {
			return StateInvalid, erroring.NewParseError("DTD value too long >0xFFFFF", p.fmtLine(), nil)
		}
		if err := p.writeVtd(TokenDtdVal, p.lastOffset, p.length1, p.depth); err != nil {
			return StateInvalid, err
		}
	} else {
		if p.length1 > maxTokenLength<<1 {
			return StateInvalid, erroring.NewParseError("DTD value too long >0xFFFFF", p.fmtLine(), nil)
		}
		if err := p.writeVtd(TokenDtdVal, p.lastOffset>>1, p.length1>>1, p.depth); err != nil {
			return StateInvalid, err
		}
	}
	if err := p.nextCharAfterWs(); err != nil {
		return StateInvalid, err
	}
	if p.currentChar != '<' {
		return StateInvalid, erroring.NewParseError(erroring.InvalidChar, p.fmtLine(), nil)
	}
	return StateLtSeen, nil
}
