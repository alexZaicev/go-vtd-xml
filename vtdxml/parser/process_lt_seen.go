package parser

import "github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"

func (p *VtdParser) processLtSeen() (State, error) {
	p.lastOffset = p.offset
	if err := p.nextChar(); err != nil {
		return StateInvalid, erroring.NewInternalError("failed to parse LT seen", err)
	}
	if p.xmlChar.IsNameStartChar(p.currentChar) {
		p.depth++
		return StateTagStart, nil
	} else {
		switch p.currentChar {
		case '/':
			return StateTagEnd, nil
		case '!':
			return p.processExSeen()
		case '?':
			return p.processQmSeen()
		default:
			return StateInvalid, erroring.NewParseError("invalid character after <", p.fmtLine(), nil)
		}
	}
}

func (p *VtdParser) processExSeen() (State, error) {
	ch, err := p.getChar()
	if err != nil {
		return StateInvalid, erroring.NewInternalError("failed to parse EX seen", err)
	}
	switch ch {
	case '-':
		{
			if p.skipChar('-') {
				p.lastOffset = p.offset
				return StateStartComment, nil
			} else {
				return StateInvalid, erroring.NewParseError("invalid char sequence to start a comment",
					p.fmtLine(), nil)
			}
		}
	case '[':
		if err := p.validateSeq(cdata); err != nil {
			return StateInvalid, err
		}
		p.lastOffset = p.offset
		return StateCdata, nil
	case 'D':
		if err := p.validateSeq(docType); err != nil {
			return StateInvalid, err
		}
		p.lastOffset = p.offset
		return StateDocType, nil
	default:
		return StateInvalid, erroring.NewParseError("unrecognized char after <!", p.fmtLine(), nil)
	}
}

func (p *VtdParser) processQmSeen() (State, error) {
	p.lastOffset = p.offset
	ch, err := p.getChar()
	if err != nil {
		return StateInvalid, erroring.NewInternalError("failed to parse QM seen", err)
	}
	if p.xmlChar.IsNameStartChar(ch) {
		if (ch == 'x' || ch == 'X') &&
			(p.skipChar('m') || p.skipChar('M')) &&
			(p.skipChar('l') || p.skipChar('L')) {
			ch, err = p.getChar()
			if err != nil {
				return StateInvalid, err
			}
			if ch == '?' || p.xmlChar.IsSpaceChar(ch) {
				return StateInvalid, erroring.NewParseError("[xX][mM][lL] not a valid PI target name",
					p.fmtLine(), nil)
			}
			offset, err := p.getPrevOffset()
			if err != nil {
				return StateInvalid, err
			}
			p.offset = offset
		}
		return StatePiTag, nil
	}
	return StateInvalid, erroring.NewParseError("invalid first char after <?", p.fmtLine(), nil)
}
