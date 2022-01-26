package parser

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

type State int

const (
	StateLtSeen State = iota
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
	StateInvalid State = -1
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

	XMLNS1998 = "http://www.w3.org/XML/1998/namespace"
	XMLNS2000 = "http://www.w3.org/2000/xmlns"
)

// nextCharAfterWs function reads next character and updates current
// and last characters
func (p *VtdParser) nextChar() error {
	if ch, err := p.reader.GetChar(); err != nil {
		return err
	} else {
		p.lastChar = p.currentChar
		p.currentChar = ch
	}
	return nil
}

// nextCharAfterWs function reads next character after whitespace
// and updates current and last characters
func (p *VtdParser) nextCharAfterWs() error {
	for {
		if err := p.nextChar(); err != nil {
			return err
		}
		if !p.xmlChar.IsSpaceChar(p.currentChar) {
			return nil
		}
	}
}

// writeVtd function writes into VTD buffer
func (p *VtdParser) writeVtd(tokenType Token, offset, length, depth int) error {
	offset64, length64, depth64 := int64(offset), int64(length), int64(depth)
	a := int64(tokenType << 28)
	b := (a | ((depth64 & 0xff) << 20) | length64) << 32
	return p.vtdBuffer.Append(b | offset64)
}

// writeVtdL5 function writes into VTD buffer and location cache
func (p *VtdParser) writeVtdL5(tokenType Token, offset, length, depth int) error {
	if err := p.writeVtd(tokenType, offset, length, depth); err != nil {
		return err
	}
	switch depth {
	case 0:
		{
			// TODO check if this can be moved to default
			p.rootIndex = p.vtdBuffer.GetSize() - 1
			break
		}
	case 1:
		{
			if p.lastDepth == 1 {
				if err := p.l1Buffer.Append(int64((p.lastL1Index << 32) | 0xFFFF)); err != nil {
					return err
				}
			} else if p.lastDepth == 2 {
				if err := p.l2Buffer.Append(int64((p.lastL2Index << 32) | 0xFFFF)); err != nil {
					return err
				}
			} else if p.lastDepth == 3 {
				if err := p.l3Buffer.Append(int64((p.lastL3Index << 32) | 0xFFFF)); err != nil {
					return err
				}
			} else if p.lastDepth == 4 {
				if err := p.l4Buffer.Append(int64((p.lastL4Index << 32) | 0xFFFF)); err != nil {
					return err
				}
			}
			p.lastL1Index = p.vtdBuffer.GetSize() - 1
			p.lastDepth = 1
			break
		}
	case 2:
		{
			if p.lastDepth == 1 {
				if err := p.l1Buffer.Append(int64((p.lastL1Index << 32) + p.l2Buffer.GetSize())); err != nil {
					return err
				}
			} else if p.lastDepth == 2 {
				if err := p.l2Buffer.Append(int64((p.lastL2Index << 32) | 0xFFFF)); err != nil {
					return err
				}
			} else if p.lastDepth == 3 {
				if err := p.l3Buffer.Append(int64((p.lastL3Index << 32) | 0xFFFF)); err != nil {
					return err
				}
			} else if p.lastDepth == 4 {
				if err := p.l4Buffer.Append(int64((p.lastL4Index << 32) | 0xFFFF)); err != nil {
					return err
				}
			}
			p.lastL2Index = p.vtdBuffer.GetSize() - 1
			p.lastDepth = 2
			break
		}
	case 3:
		{
			if p.lastDepth == 2 {
				if err := p.l2Buffer.Append(int64((p.lastL2Index << 32) + p.l3Buffer.GetSize())); err != nil {
					return err
				}
			} else if p.lastDepth == 3 {
				if err := p.l3Buffer.Append(int64((p.lastL3Index << 32) | 0xFFFF)); err != nil {
					return err
				}
			} else if p.lastDepth == 4 {
				if err := p.l4Buffer.Append(int64((p.lastL4Index << 32) | 0xFFFF)); err != nil {
					return err
				}
			}
			p.lastL3Index = p.vtdBuffer.GetSize() - 1
			p.lastDepth = 3
			break
		}
	case 4:
		{
			if p.lastDepth == 3 {
				if err := p.l3Buffer.Append(int64((p.lastL3Index << 32) + p.l4Buffer.GetSize())); err != nil {
					return err
				}
			} else if p.lastDepth == 4 {
				if err := p.l4Buffer.Append(int64((p.lastL4Index << 32) | 0xFFFF)); err != nil {
					return err
				}
			}
			p.lastL4Index = p.vtdBuffer.GetSize() - 1
			p.lastDepth = 4
			break
		}
	case 5:
		{
			if err := p.l5Buffer.Append(int64(p.vtdBuffer.GetSize() - 1)); err != nil {
				return err
			}
			if p.lastDepth == 4 {
				if err := p.l4Buffer.Append(int64((p.lastL4Index << 32) + p.l5Buffer.GetSize() - 1)); err != nil {
					return err
				}
			}
			p.lastDepth = 5
			break
		}
	}
	return nil
}

// writeVtdText function writes token text into VTD buffer
func (p *VtdParser) writeVtdText(tokenType Token, offset, length, depth int) error {
	if length > maxTokenLength {
		var j, rOffset int
		for j = length; j > maxTokenLength; j -= maxTokenLength {
			if err := p.writeVtd(tokenType, rOffset, maxTokenLength, depth); err != nil {
				return err
			}
			rOffset += maxTokenLength
		}
		return p.writeVtd(tokenType, rOffset, j, depth)
	} else {
		return p.writeVtd(tokenType, offset, length, depth)
	}
}

// fmtLine function format error message with current line number and
// character offset
func (p *VtdParser) fmtLine() string {
	return p.fmtCustomLine(p.offset)
}

// fmtLine function format error message with custom line number and
// character offset
func (p *VtdParser) fmtCustomLine(offset int) string {
	so := p.docOffset
	lineNumber, lineOffset := 0, 0

	if p.encoding < FormatUtf16BE {
		for so <= offset-1 {
			if p.xmlDoc[so] == '\n' {
				lineNumber++
				lineOffset = so
			}
			so++
		}
		lineOffset = offset - lineOffset
	} else if p.encoding == FormatUtf16BE {
		for so <= offset-2 {
			if p.xmlDoc[so+1] == '\n' && p.xmlDoc[so] == 0 {
				lineNumber++
				lineOffset = so
			}
			so += 2
		}
		lineOffset = (offset - lineOffset) >> 1
	} else {
		for so <= offset-2 {
			if p.xmlDoc[so] == '\n' && p.xmlDoc[so+1] == 0 {
				lineNumber++
				lineOffset = so
			}
			so += 2
		}
		lineOffset = (offset - lineOffset) >> 1
	}
	return fmt.Sprintf("\nLine number: %d Offset: %d", lineNumber+1, lineOffset+1)
}

// getPrevOffset function returns previous offset depending on XML document encoding
func (p *VtdParser) getPrevOffset() (int, error) {
	prevOffset := p.offset
	switch p.encoding {
	case FormatUtf8:
		for p.xmlDoc[prevOffset] < 0 && p.xmlDoc[prevOffset]&byte(0xc0) == byte(0x80) {
			prevOffset--
		}
		return prevOffset, nil
	case FormatAscii, FormatIso88591, FormatIso88592, FormatIso88593, FormatIso88594,
		FormatIso88595, FormatIso88596, FormatIso88597, FormatIso88598, FormatIso88599,
		FormatIso885910, FormatIso885911, FormatIso885912, FormatIso885913, FormatIso885914,
		FormatIso885915, FormatIso885916, FormatWin1250, FormatWin1251, FormatWin1252,
		FormatWin1253, FormatWin1254, FormatWin1255, FormatWin1256, FormatWin1257, FormatWin1258:
		return p.offset - 1, nil
	case FormatUtf16LE, FormatUtf16BE:
		temp := uint32((p.xmlDoc[p.offset]&(0xff))<<8 | (p.xmlDoc[p.offset+1] & 0xff))
		if temp < 0xd800 || temp > 0xdfff {
			return p.offset - 2, nil
		}
		return p.offset - 4, nil
	}
	return -1, erroring.NewInternalError("unsupported encoding", nil)
}

func (p *VtdParser) decideEncoding() error {
	if int32(p.xmlDoc[p.offset]) == -2 {
		p.increment = 2
		if int32(p.xmlDoc[p.offset+1]) == -1 {
			p.offset += 2
			p.encoding = FormatUtf16BE
			p.bomDetected = true
			// g.reader = reader.NewUtf16BeReader()
		} else {
			return erroring.NewEncodingError("should be 0xFF 0xFE")
		}
	} else if int32(p.xmlDoc[p.offset]) == -1 {
		p.increment = 2
		if int32(p.xmlDoc[p.offset+1]) == -2 {
			p.offset += 2
			p.encoding = FormatUtf16LE
			p.bomDetected = true
			// g.reader = reader.NewUtf16LeReader()
		} else {
			return erroring.NewEncodingError("not UTF-16LE")
		}
	} else if int32(p.xmlDoc[p.offset]) == 0 {
		if int32(p.xmlDoc[p.offset+1]) == 0x3c &&
			int32(p.xmlDoc[p.offset+2]) == 0 &&
			int32(p.xmlDoc[p.offset+3]) == 0x3f {
			p.encoding = FormatUtf16BE
			p.increment = 2
			// g.reader = reader.NewUtf16BeReader()
		} else {
			return erroring.NewEncodingError("not UTF-16BE")
		}
	} else if int32(p.xmlDoc[p.offset]) == 0x3c {
		if int32(p.xmlDoc[p.offset+1]) == 0x3c &&
			int32(p.xmlDoc[p.offset+2]) == 0 &&
			int32(p.xmlDoc[p.offset+3]) == 0x3f {
			p.encoding = FormatUtf16LE
			p.increment = 2
			// g.reader = reader.NewUtf16LeReader()
		} else {
			return erroring.NewEncodingError("not UTF-16LE")
		}
	} else if int32(p.xmlDoc[p.offset]) == -17 {
		if int32(p.xmlDoc[p.offset+1]) == -69 &&
			int32(p.xmlDoc[p.offset+2]) == -65 {
			p.offset += 3
			p.mustUtf8 = true
		} else {
			return erroring.NewEncodingError("not UTF-8")
		}
	}

	if p.encoding < FormatUtf16BE {
		if p.nsAware {
			if (p.offset + p.docLength) >= 1<<30 {
				return erroring.NewInternalError("file size too big >= 1GB", nil)
			}
		} else {
			if (p.offset + p.docLength) >= 1<<31 {
				return erroring.NewInternalError("file size too big >= 2GB", nil)
			}
		}
	} else {
		if (p.offset + p.docLength) >= 1<<31 {
			return erroring.NewInternalError("file size too big >= 2GB", nil)
		}
	}

	if p.encoding >= FormatUtf16BE {
		p.singleByteEncoding = false
	}
	return nil
}

// checkXmlPrefix functions checks XML character sequence in XML document
// byte array starting from offset. Length is passed to validate correct length
// of the expected sequence
func (p *VtdParser) checkXmlPrefix(offset, length int, checkLength bool) bool {
	var valid bool
	if p.encoding < FormatUtf16BE {
		valid = p.xmlDoc[offset] == 'x' &&
			p.xmlDoc[offset+1] == 'm' &&
			p.xmlDoc[offset+2] == 'l'
	} else if p.encoding == FormatUtf16BE {
		valid = p.xmlDoc[offset] == 0 &&
			p.xmlDoc[offset+1] == 'x' &&
			p.xmlDoc[offset+2] == 0 &&
			p.xmlDoc[offset+3] == 'm' &&
			p.xmlDoc[offset+4] == 0 &&
			p.xmlDoc[offset+5] == 'l'
	} else {
		valid = p.xmlDoc[offset] == 'x' &&
			p.xmlDoc[offset+1] == 0 &&
			p.xmlDoc[offset+2] == 'm' &&
			p.xmlDoc[offset+3] == 0 &&
			p.xmlDoc[offset+4] == 'l' &&
			p.xmlDoc[offset+5] == 0
	}

	if valid && checkLength {
		valid = (p.encoding < FormatUtf16BE && length == 4) || length == 8
	}
	return valid
}

// checkXmlnsPrefix functions checks XMLNS character sequence in XML document
// byte array starting from offset. Length is passed to validate correct length
// of the expected sequence
func (p *VtdParser) checkXmlnsPrefix(offset, length int, checkLength bool) bool {
	var valid bool
	if p.encoding < FormatUtf16BE {
		valid = p.xmlDoc[offset] == 'x' &&
			p.xmlDoc[offset+1] == 'm' &&
			p.xmlDoc[offset+2] == 'l' &&
			p.xmlDoc[offset+3] == 'n' &&
			p.xmlDoc[offset+4] == 's'
	} else if p.encoding == FormatUtf16BE {
		valid = p.xmlDoc[offset] == 0 &&
			p.xmlDoc[offset+1] == 'x' &&
			p.xmlDoc[offset+2] == 0 &&
			p.xmlDoc[offset+3] == 'm' &&
			p.xmlDoc[offset+4] == 0 &&
			p.xmlDoc[offset+5] == 'l' &&
			p.xmlDoc[offset+6] == 0 &&
			p.xmlDoc[offset+7] == 'n' &&
			p.xmlDoc[offset+8] == 0 &&
			p.xmlDoc[offset+9] == 's'
	} else {
		valid = p.xmlDoc[offset] == 'x' &&
			p.xmlDoc[offset+1] == 0 &&
			p.xmlDoc[offset+2] == 'm' &&
			p.xmlDoc[offset+3] == 0 &&
			p.xmlDoc[offset+4] == 'l' &&
			p.xmlDoc[offset+5] == 0 &&
			p.xmlDoc[offset+6] == 'n' &&
			p.xmlDoc[offset+7] == 0 &&
			p.xmlDoc[offset+8] == 's' &&
			p.xmlDoc[offset+9] == 0
	}

	if valid && checkLength {
		valid = (p.encoding < FormatUtf16BE && length == 5) || length == 10
	}
	return valid
}

// qualifyElement function does basic qualification on XML element by the
// following criteria:
// (1) current element has no prefix, then look for XMLNS;
// (2) current element has prefix, then look for XMLNS:<SOMETHING>
func (p *VtdParser) qualifyElement() error {
	preLen := int((p.currentElementRecord & 0xFFFF) >> 48)
	preOs := int(p.currentElementRecord)
	for i := p.nsBuffer3.GetSize() - 1; i >= 0; i-- {
		upVal, err := p.nsBuffer3.Upper32At(i)
		if err != nil {
			return err
		}
		diff := int((upVal & 0xFFFF) - (upVal >> 16))
		if diff == preLen {
			lowVal, err := p.nsBuffer3.Lower32At(i)
			if err != nil {
				return nil
			}
			os := int(lowVal+(upVal>>16)) + p.increment
			var j int
			for ; j < preLen-p.increment; j++ {
				if p.xmlDoc[os+j] != p.xmlDoc[preOs+j] {
					break
				}
			}
			if j == preLen-p.increment {
				return nil
			}
		}
	}
	if p.checkXmlPrefix(preOs, preLen, true) {
		return nil
	}
	return erroring.NewParseError("namespace qualification exception: element not qualified",
		p.fmtCustomLine(int(p.currentElementRecord)), nil)
}

// recordWhiteSpace function record whitespaces into VTD text buffer that are
// ignored by default
func (p *VtdParser) recordWhiteSpace() error {
	if p.depth > -1 {
		length := p.offset - p.increment - p.lastOffset
		if length != 0 {
			if p.singleByteEncoding {
				return p.writeVtdText(TokenCharacterData, p.lastOffset, length, p.depth)
			} else {
				return p.writeVtdText(TokenCharacterData, p.lastOffset>>1, length>>1, p.depth)
			}
		}
	}
	return nil
}

// handleOtherTextChar function validates other XML text character such as &amp; or &quot;
func (p *VtdParser) handleOtherTextChar(ch uint32) error {
	switch ch {
	case '&':
		{
			ch2, err := p.entityIdentifier()
			if err != nil {
				return err
			}
			if !p.xmlChar.IsValidChar(ch2) {
				return erroring.NewParseError(erroring.InvalidCharInText, p.fmtLine(), nil)
			}
			break
		}
	case ']':
		{
			// skip all ] chars
			for p.reader.SkipChar(']') {
			}
			if p.reader.SkipChar('>') {
				return erroring.NewParseError("]]> sequence in text content", p.fmtLine(), nil)
			}
			break
		}
	default:
		return erroring.NewParseError(erroring.InvalidCharInText, p.fmtLine(), nil)
	}
	return nil
}

// entityIdentifier function validates the preceding sequence of characters that
// started with & is a valid XML entity, e.g. &amp; or &quot; and returns it
// real byte equivalent
func (p *VtdParser) entityIdentifier() (uint32, error) {

	checkSeq := func(seq string) error {
		for _, seqCh := range seq {
			ch, err := p.reader.GetChar()
			if err != nil {
				return err
			}
			if int32(ch) != seqCh {
				return erroring.NewEntityError(erroring.IllegalBuiltInEntity)
			}
		}
		return nil
	}

	ch, err := p.reader.GetChar()
	if err != nil {
		return 0, err
	}
	switch ch {
	case '#':
		{
			var value uint32
			ch, err = p.reader.GetChar()
			if err != nil {
				return 0, err
			}
			if ch == 'x' {
				for {
					ch, err = p.reader.GetChar()
					if ch >= '0' && ch <= '9' {
						value = (value << 4) + (ch - '0')
					} else if ch >= 'a' && ch <= 'f' {
						value = (value << 4) + (ch - 'a' + 10)
					} else if ch >= 'A' && ch <= 'F' {
						value = (value << 4) + (ch - 'A' + 10)
					} else if ch == ';' {
						break
					} else {
						return 0, erroring.NewEntityError("illegal char following &#x")
					}
				}
			} else {
				for {
					if ch >= '0' && ch <= '9' {
						value = value*10 + (ch - '0')
					} else if ch == ';' {
						break
					} else {
						return 0, erroring.NewEntityError("illegal char following &#x")
					}
					ch, err = p.reader.GetChar()
					if err != nil {
						return 0, err
					}
				}
			}
			if !p.xmlChar.IsValidChar(value) {
				return 0, erroring.NewEntityError(erroring.InvalidChar)
			}
			return value, nil
		}
	case 'a':
		{
			// checks that the sequence matcher &amp;
			ch2, err := p.reader.GetChar()
			if err != nil {
				return 0, err
			}
			if ch2 == 'm' {
				// checks that the sequence matcher &amp;
				if err := checkSeq("p;"); err != nil {
					return 0, err
				}
				return '&', nil
			} else if ch2 == 'p' {
				// checks that the sequence matcher &apos;
				if err := checkSeq("os;"); err != nil {
					return 0, err
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
				return 0, err
			}
			return '"', nil
		}
	case 'g', 'l':
		{
			// checks that the sequence matcher &gt; or &lt;
			if err := checkSeq("t;"); err != nil {
				return 0, err
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
