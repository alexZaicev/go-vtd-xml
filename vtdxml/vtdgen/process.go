package vtdgen

import (
	"github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"
)

func (g *VtdGen) processLtSeen() (ParserState, error) {
	ch, err := g.reader.GetChar()
	if err != nil {
		return StateInvalid, erroring.NewInternalError("failed to parse LT seen", err)
	}
	if g.xmlChar.IsNameStartChar(ch) {
		g.depth++
		return StateStartTag, nil
	} else {
		switch ch {
		case '/':
			return StateEndTag, nil
		case '!':
			return g.processExSeen()
		case '?':
			return g.processQmSeen()
		default:
			return StateInvalid, erroring.NewParseError("invalid character after <", g.formatLineNumber(), nil)
		}
	}
}

func (g *VtdGen) processExSeen() (ParserState, error) {
	ch, err := g.reader.GetChar()
	if err != nil {
		return StateInvalid, erroring.NewInternalError("failed to parse EX seen", err)
	}
	switch ch {
	case '-':
		{
			if g.reader.SkipChar('-') {
				g.tmpOffset = g.offset
				return StateStartComment, nil
			} else {
				return StateInvalid, erroring.NewParseError("invalid char sequence to start a comment",
					g.formatLineNumber(), nil)
			}
		}
	case '[':
		if err := g.processCharSequence(cdata); err != nil {
			return StateInvalid, err
		}
		g.tmpOffset = g.offset
		return StateCdata, nil
	case 'D':
		if err := g.processCharSequence(docType); err != nil {
			return StateInvalid, err
		}
		g.tmpOffset = g.offset
		return StateDocType, nil
	default:
		return StateInvalid, erroring.NewParseError("unrecognized char after <!", g.formatLineNumber(), nil)
	}
}

func (g *VtdGen) processQmSeen() (ParserState, error) {
	g.tmpOffset = g.offset
	ch, err := g.reader.GetChar()
	if err != nil {
		return StateInvalid, erroring.NewInternalError("failed to parse QM seen", err)
	}
	if g.xmlChar.IsNameStartChar(ch) {
		if (ch == 'x' || ch == 'X') &&
			(g.reader.SkipChar('m') || g.reader.SkipChar('M')) &&
			(g.reader.SkipChar('l') || g.reader.SkipChar('L')) {
			ch, err = g.reader.GetChar()
			if ch == '?' || g.xmlChar.IsSpaceChar(ch) {
				return StateInvalid, erroring.NewParseError("[xX][mM][lL] not a valid PI target name",
					g.formatLineNumber(), nil)
			}
			offset, err := g.getPrevOffset()
			if err != nil {
				return StateInvalid, err
			}
			g.offset = offset
		}
		return StatePiTag, nil
	}
	return StateInvalid, erroring.NewParseError("invalid first char after <?", g.formatLineNumber(), nil)
}

func (g *VtdGen) processDocType() (ParserState, error) {
	z := 1
	for {
		ch, err := g.reader.GetChar()
		if err != nil {
			return StateInvalid, err
		}
		if !g.xmlChar.IsValidChar(ch) {
			return StateInvalid, erroring.NewParseError("invalid char in DOCTYPE", g.formatLineNumber(), nil)
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

	g.length1 = g.offset - g.tmpOffset - g.increment

	if g.singleByteEncoding {
		if g.length1 > maxTokenLength {
			return StateInvalid, erroring.NewParseError("DTD value too long >0xFFFFF", g.formatLineNumber(), nil)
		}
		if err := g.writeVtd(TokenDtdVal, g.tmpOffset, g.length1, g.depth); err != nil {
			return StateInvalid, err
		}
	} else {
		if g.length1 > maxTokenLength<<1 {
			return StateInvalid, erroring.NewParseError("DTD value too long >0xFFFFF", g.formatLineNumber(), nil)
		}
		if err := g.writeVtd(TokenDtdVal, g.tmpOffset>>1, g.length1>>1, g.depth); err != nil {
			return StateInvalid, err
		}
	}
	if ch, err := g.getCharAfterS(); err != nil {
		return StateInvalid, err
	} else {
		g.ch = ch
		if ch != '<' {
			return StateInvalid, erroring.NewParseError(erroring.InvalidChar, g.formatLineNumber(), nil)
		}
	}
	return StateLtSeen, nil
}

func (g *VtdGen) processStartTag() (ParserState, error) {
	for {
		ch, err := g.reader.GetChar()
		if err != nil {
			return StateInvalid, err
		}
		if !g.xmlChar.IsNameChar(ch) {
			break
		}
		if ch == ':' {
			g.length2 = g.offset - g.tmpOffset - g.increment
			if g.nsAware && g.checkXmlnsPrefix(g.tmpOffset, g.length2) {
				return StateInvalid, erroring.NewParseError("XMLNS cannot be element prefix", g.formatLineNumber(), err)
			}
		}
	}

	g.length1 = g.offset - g.tmpOffset - g.increment
	if g.depth > maxDepth {
		return StateInvalid, erroring.NewParseError(erroring.MaximumDepthExceeded, g.formatLineNumber(), nil)
	}

	x := (g.length1 << 32) + g.tmpOffset
	g.tagStack[g.depth] = int64(x)

	if g.depth > g.VtdDepth {
		g.VtdDepth = g.depth
	}

	if g.singleByteEncoding {
		if g.length2 > maxPrefixLength || g.length1 > maxQnameLength {
			return StateInvalid, erroring.NewParseError(erroring.TagPrefixQnameTooLong, g.formatLineNumber(), nil)
		}
		if g.shallowDepth {
			if err := g.writeVtd(TokenStartingTag, g.tmpOffset, (g.length2<<1)|g.length1, g.depth); err != nil {
				return StateInvalid, err
			}
		} else {
			if err := g.writeVtdL5(TokenStartingTag, g.tmpOffset, (g.length2<<1)|g.length1, g.depth); err != nil {
				return StateInvalid, err
			}
		}
	} else {
		if g.length2 > (maxPrefixLength<<1) || g.length1 > (maxQnameLength<<1) {
			return StateInvalid, erroring.NewParseError(erroring.TagPrefixQnameTooLong, g.formatLineNumber(), nil)
		}
		if g.shallowDepth {
			if err := g.writeVtd(TokenStartingTag, g.tmpOffset>>1, (g.length2<<10)|(g.length1>>1),
				g.depth); err != nil {
				return StateInvalid, err
			}
		} else {
			if err := g.writeVtdL5(TokenStartingTag, g.tmpOffset>>1, (g.length2<<10)|(g.length1>>1),
				g.depth); err != nil {
				return StateInvalid, err
			}
		}
	}

	if g.nsAware {
		if g.length2 != 0 {
			g.length2 += g.increment
			g.currentElementRecord = int64(((g.length2<<16)|g.length1)<<32) | int64(g.tmpOffset)
		} else {
			g.currentElementRecord = 0
		}

		if g.depth <= g.nsBuffer1.GetSize()-1 {
			g.nsBuffer1.SetSize(g.depth)
			t, err := g.nsBuffer1.IntAt(g.depth - 1)
			if err != nil {
				return StateInvalid, err
			}
			g.nsBuffer2.SetSize(int(t + 1))
			g.nsBuffer3.SetSize(int(t + 1))
		}
	}
	g.length2 = 0
	if g.xmlChar.IsSpaceChar(g.ch) {
		ch, err := g.getCharAfterS()
		if err != nil {
			return StateInvalid, err
		}
		g.ch = ch
		if g.xmlChar.IsNameStartChar(g.ch) {
			offset, err := g.getPrevOffset()
			if err != nil {
				return StateInvalid, err
			}
			g.tmpOffset = offset
			return StateAttrName, nil
		}
	}
	g.helper = true
	if g.ch == '/' {
		g.depth--
		g.helper = false
		ch, err := g.reader.GetChar()
		if err != nil {
			return StateInvalid, err
		}
		g.ch = ch
	}
	if g.ch == '>' {
		if g.nsAware {
			if err := g.nsBuffer1.Append(int32(g.nsBuffer3.GetSize() - 1)); err != nil {
				return StateInvalid, err
			}
			if g.currentElementRecord != 0 {
				if err := g.qualifyElement(); err != nil {
					return StateInvalid, err
				}
			}
		}
		return g.processElementTail()
	}
	return StateInvalid, erroring.NewParseError("invalid char in starting tag", g.formatLineNumber(), nil)
}

func (g *VtdGen) processElementTail() (ParserState, error) {
	if g.depth != -1 {
		g.tmpOffset = g.offset
		ch, err := g.getCharAfterS()
		if err != nil {
			return StateInvalid, err
		}
		g.ch = ch
		if ch == '<' {
			if g.ws {
				if err := g.recordWhiteSpace(); err != nil {
					return StateInvalid, err
				}
			}
			if g.reader.SkipChar('/') {
				if g.helper {
					g.length1 = g.offset - g.tmpOffset - (g.increment << 1)
					if g.singleByteEncoding {
						if err := g.writeVtdText(TokenCharacterData, g.tmpOffset, g.length1, g.depth); err != nil {
							return StateInvalid, err
						}
					} else {
						if err := g.writeVtdText(TokenCharacterData, g.tmpOffset>>1, g.length1>>1,
							g.depth); err != nil {
							return StateInvalid, err
						}
					}
				}
				return StateEndTag, nil
			}
			return StateLtSeen, nil
		} else if g.xmlChar.IsContentChar(ch) {
			return StateText, nil
		} else {
			if err := g.handleOtherTextChar(ch); err != nil {
				return StateInvalid, err
			}
			return StateText, nil
		}
	}
	return StateDocEnd, nil
}
