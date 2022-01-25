package vtdgen

import (
	"fmt"

	"github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"
)

type Token int64

const (
	TokenStartingTag Token = iota
	TokenEndingTag
	TokenAttrName
	TokenAttrNs
	TokenAttrVal
	TokenCharacterData
	TokenComment
	TokenPiName
	TokenPiVal
	TokenDecAttrName
	TokenDecAttrVal
	TokenCdataVal
	TokenDtdVal
	TokenDocument
)

type ParserState int

const (
	StateLtSeen ParserState = iota
	StateStartTag
	StateEndTag
	StateAttrName
	StateAttrVal
	StateText
	StateDocStart
	StateDocEnd
	StatePiTag
	StatePiVal
	StateDecAttrName
	StateStartComment
	StateEndComment
	StateCdata
	StateDocType
	StateEndPi
	StateInvalid ParserState = -1
)

type FormatEncoding int

const (
	FormatAscii FormatEncoding = iota
	FormatIso88591
	FormatIso88592
	FormatIso88593
	FormatIso88594
	FormatIso88595
	FormatIso88596
	FormatIso88597
	FormatIso88598
	FormatIso88599
	FormatIso885910
	FormatIso885911
	FormatIso885912
	FormatIso885913
	FormatIso885914
	FormatIso885915
	FormatIso885916
	FormatUtf16BE
	FormatUtf16LE
	FormatUtf8
	FormatWin1250
	FormatWin1251
	FormatWin1252
	FormatWin1253
	FormatWin1254
	FormatWin1255
	FormatWin1256
	FormatWin1257
	FormatWin1258
)

const (
	cdata   = "CDATA["
	docType = "DOCTYPE"

	maxTokenLength  = (1 << 20) - 1
	maxDepth        = 254
	maxPrefixLength = (1 << 9) - 1
	maxQnameLength  = (1 << 11) - 1
)

// writeVtd function writes into VTD buffer
func (g *VtdGen) writeVtd(tokenType Token, offset, length, depth int) error {
	offset64, length64, depth64 := int64(offset), int64(length), int64(depth)
	a := int64(tokenType << 28)
	b := (a | ((depth64 & 0xff) << 20) | length64) << 32
	return g.vtdBuffer.Append(b | offset64)
}

// writeVtdL5 function writes into VTD buffer and location cache
func (g *VtdGen) writeVtdL5(tokenType Token, offset, length, depth int) error {
	if err := g.writeVtd(tokenType, offset, length, depth); err != nil {
		return err
	}
	switch depth {
	case 0:
		{
			// TODO check if this can be moved to default
			g.rootIndex = g.vtdBuffer.GetSize() - 1
			break
		}
	case 1:
		{
			if g.lastDepth == 1 {
				if err := g.l1Buffer.Append(int64((g.lastL1Index << 32) | 0xFFFF)); err != nil {
					return err
				}
			} else if g.lastDepth == 2 {
				if err := g.l2Buffer.Append(int64((g.lastL2Index << 32) | 0xFFFF)); err != nil {
					return err
				}
			} else if g.lastDepth == 3 {
				if err := g.l3Buffer.Append(int64((g.lastL3Index << 32) | 0xFFFF)); err != nil {
					return err
				}
			} else if g.lastDepth == 4 {
				if err := g.l4Buffer.Append(int64((g.lastL4Index << 32) | 0xFFFF)); err != nil {
					return err
				}
			}
			g.lastL1Index = g.vtdBuffer.GetSize() - 1
			g.lastDepth = 1
			break
		}
	case 2:
		{
			if g.lastDepth == 1 {
				if err := g.l1Buffer.Append(int64((g.lastL1Index << 32) + g.l2Buffer.GetSize())); err != nil {
					return err
				}
			} else if g.lastDepth == 2 {
				if err := g.l2Buffer.Append(int64((g.lastL2Index << 32) | 0xFFFF)); err != nil {
					return err
				}
			} else if g.lastDepth == 3 {
				if err := g.l3Buffer.Append(int64((g.lastL3Index << 32) | 0xFFFF)); err != nil {
					return err
				}
			} else if g.lastDepth == 4 {
				if err := g.l4Buffer.Append(int64((g.lastL4Index << 32) | 0xFFFF)); err != nil {
					return err
				}
			}
			g.lastL2Index = g.vtdBuffer.GetSize() - 1
			g.lastDepth = 2
			break
		}
	case 3:
		{
			if g.lastDepth == 2 {
				if err := g.l2Buffer.Append(int64((g.lastL2Index << 32) + g.l3Buffer.GetSize())); err != nil {
					return err
				}
			} else if g.lastDepth == 3 {
				if err := g.l3Buffer.Append(int64((g.lastL3Index << 32) | 0xFFFF)); err != nil {
					return err
				}
			} else if g.lastDepth == 4 {
				if err := g.l4Buffer.Append(int64((g.lastL4Index << 32) | 0xFFFF)); err != nil {
					return err
				}
			}
			g.lastL3Index = g.vtdBuffer.GetSize() - 1
			g.lastDepth = 3
			break
		}
	case 4:
		{
			if g.lastDepth == 3 {
				if err := g.l3Buffer.Append(int64((g.lastL3Index << 32) + g.l4Buffer.GetSize())); err != nil {
					return err
				}
			} else if g.lastDepth == 4 {
				if err := g.l4Buffer.Append(int64((g.lastL4Index << 32) | 0xFFFF)); err != nil {
					return err
				}
			}
			g.lastL4Index = g.vtdBuffer.GetSize() - 1
			g.lastDepth = 4
			break
		}
	case 5:
		{
			if err := g.l5Buffer.Append(int64(g.vtdBuffer.GetSize() - 1)); err != nil {
				return err
			}
			if g.lastDepth == 4 {
				if err := g.l4Buffer.Append(int64((g.lastL4Index << 32) + g.l5Buffer.GetSize() - 1)); err != nil {
					return err
				}
			}
			g.lastDepth = 5
			break
		}
	}
	return nil
}

// writeVtdText function writes token text into VTD buffer
func (g *VtdGen) writeVtdText(tokenType Token, offset, length, depth int) error {
	if length > maxTokenLength {
		var j, rOffset int
		for j = length; j > maxTokenLength; j -= maxTokenLength {
			if err := g.writeVtd(tokenType, rOffset, maxTokenLength, depth); err != nil {
				return err
			}
			rOffset += maxTokenLength
		}
		return g.writeVtd(tokenType, rOffset, j, depth)
	} else {
		return g.writeVtd(tokenType, offset, length, depth)
	}
}

func (g *VtdGen) formatLineNumber() string {
	return g.formatCustomLineNumber(g.offset)
}

func (g *VtdGen) formatCustomLineNumber(offset int) string {
	so := g.docOffset
	lineNumber, lineOffset := 0, 0

	if g.encoding < FormatUtf16BE {
		for so <= offset-1 {
			if g.xmlDoc[so] == '\n' {
				lineNumber++
				lineOffset = so
			}
			so++
		}
		lineOffset = offset - lineOffset
	} else if g.encoding == FormatUtf16BE {
		for so <= offset-2 {
			if g.xmlDoc[so+1] == '\n' && g.xmlDoc[so] == 0 {
				lineNumber++
				lineOffset = so
			}
			so += 2
		}
		lineOffset = (offset - lineOffset) >> 1
	} else {
		for so <= offset-2 {
			if g.xmlDoc[so] == '\n' && g.xmlDoc[so+1] == 0 {
				lineNumber++
				lineOffset = so
			}
			so += 2
		}
		lineOffset = (offset - lineOffset) >> 1
	}
	return fmt.Sprintf("\nLine number: %d Offset: %d", lineNumber+1, lineOffset+1)
}

func (g *VtdGen) processCharSequence(seq string) error {
	for _, seqChar := range seq {
		c, err := g.reader.GetChar()
		if err != nil {
			return err
		}
		if c != uint32(seqChar) {
			return erroring.NewParseError(fmt.Sprintf("invalid char sequence in %s", seq),
				g.formatLineNumber(), nil)
		}
	}
	if g.depth < 0 {
		return erroring.NewParseError(fmt.Sprintf("wrong place for %s", seq),
			g.formatLineNumber(), nil)
	}
	return nil
}

func (g *VtdGen) getPrevOffset() (int, error) {
	prevOffset := g.offset
	switch g.encoding {
	case FormatUtf8:
		for g.xmlDoc[prevOffset] < 0 && g.xmlDoc[prevOffset]&byte(0xc0) == byte(0x80) {
			prevOffset--
		}
		return prevOffset, nil
	case FormatAscii, FormatIso88591, FormatIso88592, FormatIso88593, FormatIso88594,
		FormatIso88595, FormatIso88596, FormatIso88597, FormatIso88598, FormatIso88599,
		FormatIso885910, FormatIso885911, FormatIso885912, FormatIso885913, FormatIso885914,
		FormatIso885915, FormatIso885916, FormatWin1250, FormatWin1251, FormatWin1252,
		FormatWin1253, FormatWin1254, FormatWin1255, FormatWin1256, FormatWin1257, FormatWin1258:
		return g.offset - 1, nil
	case FormatUtf16LE, FormatUtf16BE:
		temp := uint32((g.xmlDoc[g.offset]&(0xff))<<8 | (g.xmlDoc[g.offset+1] & 0xff))
		if temp < 0xd800 || temp > 0xdfff {
			return g.offset - 2, nil
		}
		return g.offset - 4, nil
	}
	return -1, erroring.NewInternalError("unsupported encoding", nil)
}

func (g *VtdGen) decideEncoding() error {
	if int32(g.xmlDoc[g.offset]) == -2 {
		g.increment = 2
		if int32(g.xmlDoc[g.offset+1]) == -1 {
			g.offset += 2
			g.encoding = FormatUtf16BE
			g.bomDetected = true
			// g.reader = reader.NewUtf16BeReader()
		} else {
			return erroring.NewEncodingError("should be 0xFF 0xFE")
		}
	} else if int32(g.xmlDoc[g.offset]) == -1 {
		g.increment = 2
		if int32(g.xmlDoc[g.offset+1]) == -2 {
			g.offset += 2
			g.encoding = FormatUtf16LE
			g.bomDetected = true
			// g.reader = reader.NewUtf16LeReader()
		} else {
			return erroring.NewEncodingError("not UTF-16LE")
		}
	} else if int32(g.xmlDoc[g.offset]) == 0 {
		if int32(g.xmlDoc[g.offset+1]) == 0x3c &&
			int32(g.xmlDoc[g.offset+2]) == 0 &&
			int32(g.xmlDoc[g.offset+3]) == 0x3f {
			g.encoding = FormatUtf16BE
			g.increment = 2
			// g.reader = reader.NewUtf16BeReader()
		} else {
			return erroring.NewEncodingError("not UTF-16BE")
		}
	} else if int32(g.xmlDoc[g.offset]) == 0x3c {
		if int32(g.xmlDoc[g.offset+1]) == 0x3c &&
			int32(g.xmlDoc[g.offset+2]) == 0 &&
			int32(g.xmlDoc[g.offset+3]) == 0x3f {
			g.encoding = FormatUtf16LE
			g.increment = 2
			// g.reader = reader.NewUtf16LeReader()
		} else {
			return erroring.NewEncodingError("not UTF-16LE")
		}
	} else if int32(g.xmlDoc[g.offset]) == -17 {
		if int32(g.xmlDoc[g.offset+1]) == -69 &&
			int32(g.xmlDoc[g.offset+2]) == -65 {
			g.offset += 3
			g.mustUtf8 = true
		} else {
			return erroring.NewEncodingError("not UTF-8")
		}
	}

	if g.encoding < FormatUtf16BE {
		if g.nsAware {
			if (g.offset + g.docLength) >= 1<<30 {
				return erroring.NewInternalError("file size too big >= 1GB", nil)
			}
		} else {
			if (g.offset + g.docLength) >= 1<<31 {
				return erroring.NewInternalError("file size too big >= 2GB", nil)
			}
		}
	} else {
		if (g.offset + g.docLength) >= 1<<31 {
			return erroring.NewInternalError("file size too big >= 2GB", nil)
		}
	}

	if g.encoding >= FormatUtf16BE {
		g.singleByteEncoding = false
	}
	return nil
}

func (g *VtdGen) getCharAfterS() (uint32, error) {
	for {
		ch, err := g.reader.GetChar()
		if err != nil {
			return 0, err
		}
		if !g.xmlChar.IsSpaceChar(ch) {
			return ch, nil
		}
	}
}

// checkXmlPrefix functions checks XML character sequence in XML document
// byte array starting from offset. Length is passed to validate correct length
// of the expected sequence
func (g *VtdGen) checkXmlPrefix(offset, length int) bool {
	if g.encoding < FormatUtf16BE {
		return length == 4 &&
			g.xmlDoc[offset] == 'x' &&
			g.xmlDoc[offset+1] == 'm' &&
			g.xmlDoc[offset+2] == 'l'
	} else if g.encoding == FormatUtf16BE {
		return length == 8 &&
			g.xmlDoc[offset] == 0 &&
			g.xmlDoc[offset+1] == 'x' &&
			g.xmlDoc[offset+2] == 0 &&
			g.xmlDoc[offset+3] == 'm' &&
			g.xmlDoc[offset+4] == 0 &&
			g.xmlDoc[offset+5] == 'l'
	} else {
		return length == 8 &&
			g.xmlDoc[offset] == 'x' &&
			g.xmlDoc[offset+1] == 0 &&
			g.xmlDoc[offset+2] == 'm' &&
			g.xmlDoc[offset+3] == 0 &&
			g.xmlDoc[offset+4] == 'l' &&
			g.xmlDoc[offset+5] == 0
	}
}

// checkXmlnsPrefix functions checks XMLNS character sequence in XML document
// byte array starting from offset. Length is passed to validate correct length
// of the expected sequence
func (g *VtdGen) checkXmlnsPrefix(offset, length int) bool {
	if g.encoding < FormatUtf16BE {
		return length == 5 &&
			g.xmlDoc[offset] == 'x' &&
			g.xmlDoc[offset+1] == 'm' &&
			g.xmlDoc[offset+2] == 'l' &&
			g.xmlDoc[offset+3] == 'n' &&
			g.xmlDoc[offset+4] == 's'
	} else if g.encoding == FormatUtf16BE {
		return length == 10 &&
			g.xmlDoc[offset] == 0 &&
			g.xmlDoc[offset+1] == 'x' &&
			g.xmlDoc[offset+2] == 0 &&
			g.xmlDoc[offset+3] == 'm' &&
			g.xmlDoc[offset+4] == 0 &&
			g.xmlDoc[offset+5] == 'l' &&
			g.xmlDoc[offset+6] == 0 &&
			g.xmlDoc[offset+7] == 'n' &&
			g.xmlDoc[offset+8] == 0 &&
			g.xmlDoc[offset+9] == 's'
	} else {
		return length == 10 &&
			g.xmlDoc[offset] == 'x' &&
			g.xmlDoc[offset+1] == 0 &&
			g.xmlDoc[offset+2] == 'm' &&
			g.xmlDoc[offset+3] == 0 &&
			g.xmlDoc[offset+4] == 'l' &&
			g.xmlDoc[offset+5] == 0 &&
			g.xmlDoc[offset+6] == 'n' &&
			g.xmlDoc[offset+7] == 0 &&
			g.xmlDoc[offset+8] == 's' &&
			g.xmlDoc[offset+9] == 0
	}
}

// qualifyElement function does basic qualification on XML element by the
// following criteria:
// (1) current element has no prefix, then look for XMLNS;
// (2) current element has prefix, then look for XMLNS:<SOMETHING>
func (g *VtdGen) qualifyElement() error {
	preLen := int((g.currentElementRecord & 0xFFFF) >> 48)
	preOs := int(g.currentElementRecord)
	for i := g.nsBuffer3.GetSize() - 1; i >= 0; i-- {
		upVal, err := g.nsBuffer3.Upper32At(i)
		if err != nil {
			return err
		}
		diff := int((upVal & 0xFFFF) - (upVal >> 16))
		if diff == preLen {
			lowVal, err := g.nsBuffer3.Lower32At(i)
			if err != nil {
				return nil
			}
			os := int(lowVal+(upVal>>16)) + g.increment
			var j int
			for ; j < preLen-g.increment; j++ {
				if g.xmlDoc[os+j] != g.xmlDoc[preOs+j] {
					break
				}
			}
			if j == preLen-g.increment {
				return nil
			}
		}
	}
	if g.checkXmlPrefix(preOs, preLen) {
		return nil
	}
	return erroring.NewParseError("namespace qualification exception: element not qualified",
		g.formatCustomLineNumber(int(g.currentElementRecord)), nil)
}

func (g *VtdGen) recordWhiteSpace() error {
	if g.depth > -1 {
		length := g.offset - g.increment - g.tmpOffset
		if length != 0 {
			if g.singleByteEncoding {
				return g.writeVtdText(TokenCharacterData, g.tmpOffset, length, g.depth)
			} else {
				return g.writeVtdText(TokenCharacterData, g.tmpOffset>>1, length>>1, g.depth)
			}

		}
	}
	return nil
}

func (g *VtdGen) handleOtherTextChar(ch uint32) error {
	switch ch {
	case '&':
		{
			ch2, err := g.entityIdentifier()
			if err != nil {
				return err
			}
			if !g.xmlChar.IsValidChar(ch2) {
				return erroring.NewParseError(erroring.InvalidCharInText, g.formatLineNumber(), nil)
			}
			break
		}
	case ']':
		{
			// skip all ] chars
			for g.reader.SkipChar(']') {
			}
			if g.reader.SkipChar('>') {
				return erroring.NewParseError("]]> sequence in text content", g.formatLineNumber(), nil)
			}
			break
		}
	default:
		return erroring.NewParseError(erroring.InvalidCharInText, g.formatLineNumber(), nil)
	}
	return nil
}

func (g *VtdGen) entityIdentifier() (uint32, error) {

	checkSeq := func(seq string) error {
		for _, seqCh := range seq {
			ch, err := g.reader.GetChar()
			if err != nil {
				return err
			}
			if int32(ch) != seqCh {
				return erroring.NewEntityError(erroring.IllegalBuiltInEntity)
			}
		}
		return nil
	}

	ch, err := g.reader.GetChar()
	if err != nil {
		return 0, err
	}
	switch ch {
	case '#':
		{
		
		}
	case 'a':
		{
			// checks that the sequence matcher &amp;
			ch2, err := g.reader.GetChar()
			if err != nil {
				return 0, err
			}
			if ch2 == 'm' {
				// checks that the sequence matcher &amp;
				if err := checkSeq("p;"); err != nil {
					return 0, nil
				}
				return '&', nil
			} else if ch2 == 'p' {
				// checks that the sequence matcher &apos;
				if err := checkSeq("os;"); err != nil {
					return 0, nil
				}
				return '\'', nil
			} else {
				return 0, erroring.NewEntityError(erroring.IllegalBuiltInEntity)
			}
		}
	case 'q':
		{
			// checks that the sequence matcher &quot;
			if err := checkSeq("uot;"); err != nil {
				return 0, nil
			}
			return '"', nil
		}
	case 'g', 'l':
		{
			// checks that the sequence matcher &gt; or &lt;
			if err := checkSeq("t;"); err != nil {
				return 0, nil
			}
			if ch == 'g' {
				return '>', nil
			} else {
				return '<', nil
			}
		}
	default:
		return 0, erroring.NewEntityError("illegal entity character")
	}
}
