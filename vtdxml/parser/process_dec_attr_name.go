package parser

import (
	"strings"

	"github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"
	"github.com/alexZaicev/go-vtd-xml/vtdxml/reader"
)

const (
	version    = "version"
	encoding   = "encoding"
	standalone = "standalone"

	ascii   = "ascii"
	usAscii = "us-ascii"
	cp125   = "cp125"
	windows = "windows-125"
	iso8859 = "iso-8859-"
	utf8    = "utf-8"
	utf16   = "utf-16"
	utf16LE = "utf-16le"
	utf16BE = "utf-16be"
)

func (p *VtdParser) processDecAttrName() (State, error) {
	if p.currentChar == uint32(version[0]) && p.skipCharSeq(version[1:]) {
		if err := p.nextCharAfterWs(); err != nil {
			return StateInvalid, err
		}
		if p.currentChar != '=' {
			return StateInvalid, erroring.NewParseError(erroring.InvalidChar, p.fmtLine(), nil)
		}
		if p.singleByteEncoding {
			if err := p.writeVtd(TokenDecAttrName, p.lastOffset-1, 7, p.depth); err != nil {
				return StateInvalid, err
			}
		} else {
			if err := p.writeVtd(TokenDecAttrName, (p.lastOffset-2)>>2, 7, p.depth); err != nil {
				return StateInvalid, err
			}
		}
	} else {
		return StateInvalid, erroring.NewParseError("declaration should be version", p.fmtLine(), nil)
	}
	if err := p.nextCharAfterWs(); err != nil {
		return StateInvalid, err
	}
	if p.currentChar != '\'' && p.currentChar != '"' {
		return StateInvalid, erroring.NewParseError("invalid char to start attribute name", p.fmtLine(), nil)
	}
	p.lastOffset = p.offset
	// support 1.0 & 1.1 versions
	if p.skipCharSeq("1.") && (p.skipChar('0') || p.skipChar('1')) {
		if p.singleByteEncoding {
			if err := p.writeVtd(TokenDecAttrVal, p.lastOffset, 3, p.depth); err != nil {
				return StateInvalid, err
			}
		} else {
			if err := p.writeVtd(TokenDecAttrVal, p.lastOffset>>1, 3, p.depth); err != nil {
				return StateInvalid, err
			}
		}
	} else {
		return StateInvalid, erroring.NewParseError("invalid version detected (supported 1.0 or 1.1)", p.fmtLine(), nil)
	}
	if !p.skipChar(p.currentChar) {
		return StateInvalid, erroring.NewParseError("version not terminated properly", p.fmtLine(), nil)
	}
	if err := p.nextChar(); err != nil {
		return StateInvalid, err
	}

	if p.xmlChar.IsSpaceChar(p.currentChar) {
		if err := p.nextCharAfterWs(); err != nil {
			return StateInvalid, err
		}
		p.lastOffset = p.offset - p.increment
		if p.currentChar == uint32(encoding[0]) {
			if !p.skipCharSeq(encoding[1:]) {
				return StateInvalid, erroring.NewParseError("declaration should be encoding", p.fmtLine(), nil)
			}
			if err := p.processDecEncodingAttr(); err != nil {
				return StateInvalid, err
			}
		}
		if p.currentChar == uint32(standalone[0]) {
			if !p.skipCharSeq(standalone[1:]) {
				return StateInvalid, erroring.NewParseError("declaration should be standalone", p.fmtLine(), nil)
			}
			if err := p.processDecStandaloneAttr(); err != nil {
				return StateInvalid, err
			}
		}
	}
	if p.currentChar == '?' && p.skipChar('>') {
		p.lastOffset = p.offset
		if err := p.nextCharAfterWs(); err != nil {
			return StateInvalid, err
		}
		if p.currentChar == '<' {
			return StateLtSeen, nil
		} else {
			return StateInvalid, erroring.NewParseError(erroring.InvalidChar, p.fmtLine(), nil)
		}
	} else {
		return StateInvalid, erroring.NewParseError("invalid termination sequence", p.fmtLine(), nil)
	}
}

func (p *VtdParser) processDecEncodingAttr() error {
	if err := p.nextCharAfterWs(); err != nil {
		return err
	}
	if p.currentChar != '=' {
		return erroring.NewParseError(erroring.InvalidChar, p.fmtLine(), nil)
	}
	if p.singleByteEncoding {
		if err := p.writeVtd(TokenDecAttrName, p.lastOffset, 8, p.depth); err != nil {
			return err
		}
	} else {
		if err := p.writeVtd(TokenDecAttrName, p.lastOffset>>1, 8, p.depth); err != nil {
			return err
		}
	}
	if err := p.nextCharAfterWs(); err != nil {
		return err
	}
	if p.currentChar != '\'' && p.currentChar != '"' {
		return erroring.NewParseError("invalid char to start attribute name", p.fmtLine(), nil)
	}
	p.lastOffset = p.offset
	if err := p.nextChar(); err != nil {
		return err
	}
	switch p.currentChar {
	case 'a', 'A':
		if err := p.checkAsciiEncoding(); err != nil {
			return err
		}
	case 'c', 'C':
		if err := p.checkCpEncoding(); err != nil {
			return err
		}
	case 'i', 'I':
		if err := p.checkIsoEncoding(); err != nil {
			return err
		}
	case 'u', 'U':
		if p.skipChar(uint32(usAscii[1])) {
			if err := p.checkUsAsciiEncoding(); err != nil {
				return err
			}
		} else {
			if err := p.checkUtfEncoding(); err != nil {
				return err
			}
		}
	case 'w', 'W':
		if err := p.checkWindowsEncoding(); err != nil {
			return err
		}
	default:
		return erroring.NewParseError("invalid document encoding", p.fmtLine(), nil)
	}
	if err := p.nextChar(); err != nil {
		return err
	}
	if p.currentChar != '\'' && p.currentChar != '"' {
		return erroring.NewParseError("invalid char to start attribute name", p.fmtLine(), nil)
	}
	if err := p.nextCharAfterWs(); err != nil {
		return err
	}
	p.lastOffset = p.offset - p.increment
	return nil
}

func (p *VtdParser) processDecStandaloneAttr() error {
	if err := p.nextCharAfterWs(); err != nil {
		return err
	}
	if p.currentChar != '=' {
		return erroring.NewParseError(erroring.InvalidChar, p.fmtLine(), nil)
	}
	if p.singleByteEncoding {
		if err := p.writeVtd(TokenDecAttrName, p.lastOffset, 10, p.depth); err != nil {
			return err
		}
	} else {
		if err := p.writeVtd(TokenDecAttrName, p.lastOffset>>1, 10, p.depth); err != nil {
			return err
		}
	}
	if err := p.nextCharAfterWs(); err != nil {
		return err
	}
	p.lastOffset = p.offset
	if p.currentChar != '\'' && p.currentChar != '"' {
		return erroring.NewParseError("invalid char to start attribute name", p.fmtLine(), nil)
	}
	if p.skipCharSeq("yes") {
		if p.singleByteEncoding {
			if err := p.writeVtd(TokenDecAttrVal, p.lastOffset, 3, p.depth); err != nil {
				return err
			}
		} else {
			if err := p.writeVtd(TokenDecAttrVal, p.lastOffset>>1, 3, p.depth); err != nil {
				return err
			}
		}
	} else if p.skipCharSeq("no") {
		if p.singleByteEncoding {
			if err := p.writeVtd(TokenDecAttrVal, p.lastOffset, 2, p.depth); err != nil {
				return err
			}
		} else {
			if err := p.writeVtd(TokenDecAttrVal, p.lastOffset>>1, 2, p.depth); err != nil {
				return err
			}
		}
	} else {
		return erroring.NewParseError("invalid value for attribute standalone (valid options are yes or no)",
			p.fmtLine(), nil)
	}
	if err := p.nextChar(); err != nil {
		return err
	}
	if p.currentChar != '\'' && p.currentChar != '"' {
		return erroring.NewParseError("invalid char to start attribute name", p.fmtLine(), nil)
	}
	if err := p.nextCharAfterWs(); err != nil {
		return err
	}
	return nil
}

func (p *VtdParser) checkAsciiEncoding() error {
	if !p.skipCharSeq(ascii[1:]) && !p.skipCharSeq(strings.ToUpper(ascii[1:])) {
		return erroring.NewParseError("invalid document encoding", p.fmtLine(), nil)
	}
	if p.encoding < FormatUtf16BE || p.encoding == FormatUtf16LE || p.mustUtf8 {
		return erroring.NewParseError("cannot switch document encoding to ASCII", p.fmtLine(), nil)
	}
	p.encoding = FormatAscii
	r, err := reader.NewAsciiReader(p.xmlDoc, p.offset, p.endOffset)
	if err != nil {
		return err
	}
	p.reader = r
	if err := p.writeVtd(TokenDecAttrVal, p.lastOffset, 5, p.depth); err != nil {
		return err
	}
	return nil
}

func (p *VtdParser) checkUsAsciiEncoding() error {
	if !p.skipCharSeq(usAscii[2:]) && !p.skipCharSeq(strings.ToUpper(usAscii[2:])) {
		return erroring.NewParseError("invalid document encoding", p.fmtLine(), nil)
	}
	if p.encoding < FormatUtf16BE || p.encoding == FormatUtf16LE || p.mustUtf8 {
		return erroring.NewParseError("cannot switch document encoding to US-ASCII", p.fmtLine(), nil)
	}
	p.encoding = FormatAscii
	r, err := reader.NewAsciiReader(p.xmlDoc, p.offset, p.endOffset)
	if err != nil {
		return err
	}
	p.reader = r
	if err := p.writeVtd(TokenDecAttrVal, p.lastOffset, 8, p.depth); err != nil {
		return err
	}
	return nil
}

func (p *VtdParser) checkIsoEncoding() error {
	if !p.skipCharSeq(iso8859[1:]) && !p.skipCharSeq(strings.ToUpper(iso8859[1:])) {
		return erroring.NewParseError("invalid document encoding", p.fmtLine(), nil)
	}
	if p.encoding < FormatUtf16BE || p.encoding == FormatUtf16LE || p.mustUtf8 {
		return erroring.NewParseError("cannot switch document encoding to ISO8859x", p.fmtLine(), nil)
	}
	if length, err := p.setIsoEncoding(); err != nil {
		return err
	} else {
		if err := p.writeVtd(TokenDecAttrVal, p.lastOffset, length, p.depth); err != nil {
			return err
		}
	}
	return nil
}

func (p *VtdParser) checkUtfEncoding() error {
	if p.skipCharSeqIgnoreCase(utf8[1:4]) {
		if p.skipChar(uint32(utf8[4])) {
			// resolve UTF-8
			if !p.singleByteEncoding {
				return erroring.NewParseError("cannot switch document encoding to UTF-8", p.fmtLine(), nil)
			}
			if err := p.writeVtd(TokenDecAttrVal, p.lastOffset, 5, p.depth); err != nil {
				return err
			}
		} else if p.skipCharSeqIgnoreCase(utf16[4:6]) {
			if p.skipCharSeqIgnoreCase(utf16LE[6:]) {
				// resolve UTF-16LE
				if p.encoding == FormatUtf16LE {
					// r, err := reader.NewUtf16LeReader(p.xmlDoc, p.offset, p.endOffset)
					// if err != nil {
					// 	return err
					// }
					// p.reader = &r
					if err := p.writeVtd(TokenDecAttrVal, p.lastOffset>>1, 8, p.depth); err != nil {
						return err
					}
				} else {
					return erroring.NewParseError("cannot switch document encoding to UTF-16LE", p.fmtLine(), nil)
				}
			} else if p.skipCharSeqIgnoreCase(utf16BE[6:]) {
				// resolve UTF-16BE
				if p.encoding == FormatUtf16BE {
					// r, err := reader.NewUtf16BeReader(p.xmlDoc, p.offset, p.endOffset)
					// if err != nil {
					// 	return err
					// }
					// p.reader = &r
					if err := p.writeVtd(TokenDecAttrVal, p.lastOffset>>1, 8, p.depth); err != nil {
						return err
					}
				} else {
					return erroring.NewParseError("cannot switch document encoding to UTF-16BE", p.fmtLine(), nil)
				}
			} else {
				// resolve UTF-16
				if !p.singleByteEncoding {
					if !p.bomDetected {
						return erroring.NewParseError("BOM not detected for UTF-16", p.fmtLine(), nil)
					}
					// TODO identify why there is not reader for UTF-16
					if err := p.writeVtd(TokenDecAttrVal, p.lastOffset>>1, 6, p.depth); err != nil {
						return err
					}
				} else {
					return erroring.NewParseError("cannot switch document encoding to UTF-16", p.fmtLine(), nil)
				}
			}
		} else {
			return erroring.NewParseError("invalid document encoding", p.fmtLine(), nil)
		}
	} else {
		return erroring.NewParseError("invalid document encoding", p.fmtLine(), nil)
	}
	return nil
}

func (p *VtdParser) checkCpEncoding() error {
	if !p.skipCharSeq(cp125[1:]) && !p.skipCharSeq(strings.ToUpper(cp125[1:])) {
		return erroring.NewParseError("invalid document encoding", p.fmtLine(), nil)
	}
	if p.encoding > FormatUtf16LE || p.mustUtf8 {
		return erroring.NewParseError("cannot switch document encoding to CP125x", p.fmtLine(), nil)
	}
	if err := p.setWinEncoding(); err != nil {
		return err
	}
	if err := p.writeVtd(TokenDecAttrVal, p.lastOffset, 6, p.depth); err != nil {
		return err
	}
	return nil
}

func (p *VtdParser) checkWindowsEncoding() error {
	if !p.skipCharSeq(windows[1:]) && !p.skipCharSeq(strings.ToUpper(windows[1:])) {
		return erroring.NewParseError("invalid document encoding", p.fmtLine(), nil)
	}
	if p.encoding > FormatUtf16LE || p.mustUtf8 {
		return erroring.NewParseError("cannot switch document encoding to WINDOWS-125x", p.fmtLine(), nil)
	}
	if err := p.setWinEncoding(); err != nil {
		return err
	}
	if err := p.writeVtd(TokenDecAttrVal, p.lastOffset, 12, p.depth); err != nil {
		return err
	}
	return nil
}

func (p *VtdParser) setWinEncoding() error {
	// if p.skipChar('0') {
	// 	p.encoding = FormatWin1250
	// 	r, err := reader.NewWin1250(p.xmlDoc, p.offset, p.endOffset)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	p.reader = &r
	// } else if p.skipChar('1') {
	// 	p.encoding = FormatWin1251
	// 	r, err := reader.NewWin1251(p.xmlDoc, p.offset, p.endOffset)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	p.reader = &r
	// } else if p.skipChar('2') {
	// 	p.encoding = FormatWin1252
	// 	r, err := reader.NewWin1252(p.xmlDoc, p.offset, p.endOffset)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	p.reader = &r
	// } else if p.skipChar('3') {
	// 	p.encoding = FormatWin1253
	// 	r, err := reader.NewWin1253(p.xmlDoc, p.offset, p.endOffset)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	p.reader = &r
	// } else if p.skipChar('4') {
	// 	p.encoding = FormatWin1254
	// 	r, err := reader.NewWin1254(p.xmlDoc, p.offset, p.endOffset)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	p.reader = &r
	// } else if p.skipChar('5') {
	// 	p.encoding = FormatWin1255
	// 	r, err := reader.NewWin1255(p.xmlDoc, p.offset, p.endOffset)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	p.reader = &r
	// } else if p.skipChar('6') {
	// 	p.encoding = FormatWin1256
	// 	r, err := reader.NewWin1256(p.xmlDoc, p.offset, p.endOffset)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	p.reader = &r
	// } else if p.skipChar('7') {
	// 	p.encoding = FormatWin1257
	// 	r, err := reader.NewWin1257(p.xmlDoc, p.offset, p.endOffset)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	p.reader = &r
	// } else if p.skipChar('8') {
	// 	p.encoding = FormatWin1258
	// 	r, err := reader.NewWin1258(p.xmlDoc, p.offset, p.endOffset)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	p.reader = &r
	// } else {
	// 	return erroring.NewParseError("invalid document encoding", p.fmtLine(), nil)
	// }
	return nil
}

func (p *VtdParser) setIsoEncoding() (int, error) {
	// TODO implement: return length end error
	return 0, nil
}
