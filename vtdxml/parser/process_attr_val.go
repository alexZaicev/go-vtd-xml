package parser

import (
	"fmt"

	"github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"
)

type NsUrlType int

const (
	DefaultNsUrl NsUrlType = iota
	NsUrl2000
	NsUrl1998
	InvalidNsUrl NsUrlType = -1

	URL1 = "2000/xmlns/"
	URL2 = "http://www.w3.org/XML/1998/namespace"
)

type Url struct {
	offset, length int
}

func (p *VtdParser) processAttrVal() (State, error) {
	ch := p.currentChar
	for {
		if err := p.nextChar(); err != nil {
			return StateInvalid, err
		}
		if p.xmlChar.IsValidChar(p.currentChar) && p.currentChar != '<' {
			if p.currentChar == ch {
				// break because attribute value finished
				break
			}
			if p.currentChar == '&' {
				if ch, err := p.entityIdentifier(); err != nil {
					return StateInvalid, err
				} else if !p.xmlChar.IsValidChar(ch) {
					return StateInvalid, erroring.NewParseError(erroring.InvalidChar, p.fmtLine(), nil)
				}
			}
		} else {
			return StateInvalid, erroring.NewParseError(erroring.InvalidChar, p.fmtLine(), nil)
		}
	}
	p.length1 = p.offset - p.lastOffset - p.increment
	if p.nsAware && p.isNs {
		if !p.defaultNs && p.length1 == 0 {
			return StateInvalid, erroring.NewParseError(erroring.NonDefaultNsEmpty, p.fmtLine(), nil)
		}
		nsUrlType, err := p.identifyNsUrl()
		if err != nil {
			return StateInvalid, err
		}
		if p.isXml && nsUrlType != NsUrl2000 {
			return StateInvalid, erroring.NewParseError(
				fmt.Sprintf("xmlns:xml cant only point to %s", XMLNS1998),
				p.fmtLine(),
				nil,
			)
		} else {
			if !p.defaultNs {
				if err := p.nsBuffer2.Append(int64(p.lastOffset<<32 | p.length1)); err != nil {
					return StateInvalid, err
				}
			}
			if nsUrlType != DefaultNsUrl {
				if nsUrlType == NsUrl1998 {
					return StateInvalid, erroring.NewParseError(
						fmt.Sprintf("namespace declation cannot point to %s", XMLNS1998),
						p.fmtLine(),
						nil,
					)
				}
				return StateInvalid, erroring.NewParseError(
					fmt.Sprintf("namespace declation cannot point to %s", XMLNS2000),
					p.fmtLine(),
					nil,
				)
			}
		}
	}
	if err := p.writeVtdWithLengthCheck(TokenAttrVal, erroring.AttrValueTooLong); err != nil {
		return StateInvalid, err
	}
	p.isXml, p.isNs = false, false
	if err := p.nextChar(); err != nil {
		return StateInvalid, err
	}
	if p.xmlChar.IsSpaceChar(p.currentChar) {
		if err := p.nextCharAfterWs(); err != nil {
			return StateInvalid, err
		}
		if p.xmlChar.IsNameStartChar(p.currentChar) {
			p.lastOffset = p.offset - p.increment
			return StateAttrName, nil
		}
	}

	p.helper = true
	if p.currentChar == '/' {
		p.depth--
		p.helper = false
		if err := p.nextChar(); err != nil {
			return StateInvalid, err
		}
	}
	if p.currentChar == '>' {
		if p.nsAware {
			if err := p.nsBuffer1.Append(int32(p.nsBuffer3.GetSize() - 1)); err != nil {
				return StateInvalid, err
			}
			if p.prefixedAttCount > 0 {
				if err := p.qualifyAttributes(); err != nil {
					return StateInvalid, err
				}
			}
			if p.prefixedAttCount > 0 {
				if err := p.checkQualifiedAttributeUniqueness(); err != nil {
					return StateInvalid, err
				}
			}
			if p.currentElementRecord != 0 {
				if err := p.qualifyElement(); err != nil {
					return StateInvalid, err
				}
			}
			p.prefixedAttCount = 0
		}
		p.attrCount = 0
		return p.processElementTail()
	}
	return StateInvalid, erroring.NewParseError(erroring.InvalidChar, p.fmtLine(), nil)
}

func (p *VtdParser) identifyNsUrl() (NsUrlType, error) {
	// TODO instead of returning int values,
	//  create a struct that would bring meaningful explanation about the functionality
	// lastOffset, length1
	g := p.lastOffset + p.length1
	offset, tmpOffset := p.lastOffset, 0
	if p.length1 < 29 || (p.increment == 2 && p.length1 < 58) {
		return DefaultNsUrl, nil
	}

	for i := 0; i < 18 && offset < g; i++ {
		ch, err := p.getCharResolved(offset)
		if err != nil {
			return InvalidNsUrl, err
		}
		if int64(URL2[i]) != ch {
			return DefaultNsUrl, nil
		}
		offset += int(ch) >> 32
	}
	tmpOffset = offset
	for i := 0; i < 11 && offset < g; i++ {
		ch, err := p.getCharResolved(offset)
		if err != nil {
			return InvalidNsUrl, err
		}
		if int64(URL1[i]) != ch {
			return DefaultNsUrl, nil
		}
		offset += int(ch) >> 32
	}
	if offset == g {
		return NsUrl1998, nil
	}
	offset = tmpOffset
	for i := 18; i < 36 && offset < g; i++ {
		ch, err := p.getCharResolved(offset)
		if err != nil {
			return InvalidNsUrl, err
		}
		if int64(URL2[i]) != ch {
			return DefaultNsUrl, nil
		}
		offset += int(ch) >> 32
	}
	if offset == g {
		return NsUrl2000, nil
	}
	return DefaultNsUrl, nil
}

func (p *VtdParser) qualifyAttributes() error {
	nsBuffer3Count := p.nsBuffer3.GetSize() - 1
	for j := 0; j < p.prefixedAttCount; j++ {
		preLen := int(int32((p.prefixedAttrNameSlice[j] & 0xFFFF0000) >> 16))
		preOs := int(int32(p.prefixedAttrNameSlice[j] >> 32))

		i := nsBuffer3Count
		for ; i >= 0; i-- {
			upper, err := p.nsBuffer3.Upper32At(i)
			if err != nil {
				return err
			}
			if (int(upper)&0xFFFF)-(int(upper)>>16) == preLen+p.increment {
				// doing byte comparison here
				lower, err := p.nsBuffer3.Lower32At(i)
				if err != nil {
					return err
				}
				offset := int(lower) + (int(upper) >> 16) + p.increment
				var k int
				for ; k < preLen; k++ {
					if p.xmlDoc[offset+k] != p.xmlDoc[preOs+k] {
						break
					}
				}
				if k == preLen {
					// match found break
					break
				}
			}
		}
		if i < 0 {
			return erroring.NewParseError("prefixed attribute not qualified", p.fmtLine(), nil)
		} else {
			p.prefixUrlSlice = append(p.prefixUrlSlice, i)
		}
	}
	return nil
}

func (p *VtdParser) checkQualifiedAttributeUniqueness() error {
	for i := 0; i < p.prefixedAttCount; i++ {
		preLen := int(int32((p.prefixedAttrNameSlice[i] & 0xFFFF0000) >> 16))
		postLen := int(int32(p.prefixedAttrNameSlice[i]&0xFFFF)) - preLen - p.increment

		offset := int(int32(p.prefixedAttrNameSlice[i]>>32)) + preLen + p.increment

		urlLen, err := p.nsBuffer2.Lower32At(p.prefixUrlSlice[i])
		if err != nil {
			return err
		}
		urlOffset, err := p.nsBuffer2.Upper32At(p.prefixUrlSlice[i])
		if err != nil {
			return err
		}
		urlA := Url{int(urlOffset), int(urlLen)}

		for j := i + 1; j < p.prefixedAttCount; j++ {
			preLen2 := int(int32((p.prefixedAttrNameSlice[j] & 0xFFFF0000) >> 16))
			postLen2 := int(int32(p.prefixedAttrNameSlice[j]&0xFFFF)) - preLen2 - p.increment

			offset2 := int(int32(p.prefixedAttrNameSlice[j]>>32)) + preLen2 + p.increment

			if postLen == postLen2 {
				var k int
				for ; k < postLen; k++ {
					if p.xmlDoc[offset+k] != p.xmlDoc[offset2+k] {
						break
					}
				}
				if k == postLen {
					urlLen2, err := p.nsBuffer2.Lower32At(p.prefixUrlSlice[j])
					if err != nil {
						return err
					}
					urlOffset2, err := p.nsBuffer2.Upper32At(p.prefixUrlSlice[j])
					if err != nil {
						return err
					}
					urlB := Url{int(urlOffset2), int(urlLen2)}

					match, err := p.matchUrl(urlA, urlB)
					if err != nil {
						return err
					}
					if match {
						return erroring.NewParseError("qualified attribute names collide",
							p.fmtCustomLine(offset2),
							nil)
					}
				}
			}
		}
	}
	return nil
}

func (p *VtdParser) getCharResolved(offset int) (int64, error) {
	checkSeq := func(seq string, offset, inc int) error {
		var i int
		for _, seqCh := range seq {
			ch, err := p.getCharUnit(offset + i)
			if err != nil {
				return err
			}
			if ch != seqCh {
				return erroring.NewEntityError(erroring.IllegalBuiltInEntity)
			}
			i += inc
		}
		return nil
	}

	var val int32
	inc := 2 << (p.increment - 1)
	ch64, err := p.reader.GetLongChar(int32(offset))
	if err != nil {
		return 0, err
	}
	ch := int32(ch64)
	if ch != '&' {
		return int64(ch64), nil
	}
	offset += p.increment
	ch, err = p.getCharUnit(offset)
	if err != nil {
		return 0, err
	}
	offset += p.increment

	switch ch {
	case '#':
		{
			ch, err = p.getCharUnit(offset)
			if err != nil {
				return 0, err
			}
			if ch == 'x' {
				for {
					offset += p.increment
					inc += p.increment
					ch, err = p.getCharUnit(offset)
					if err != nil {
						return 0, err
					}
					if ch >= '0' && ch <= '9' {
						val = (val << 4) + (ch - '0')
					} else if ch >= 'a' && ch <= 'f' {
						val = (val << 4) + (ch - 'a' + 10)
					} else if ch >= 'A' && ch <= 'F' {
						val = (val << 4) + (ch - 'A' + 10)
					} else if ch == ';' {
						inc += p.increment
						break
					}
				}
			} else {
				for {
					ch, err = p.getCharUnit(offset)
					if err != nil {
						return 0, err
					}
					offset += p.increment
					inc += p.increment
					if ch >= '0' && ch <= '9' {
						val = val*10 + (ch - '0')
					} else if ch == ';' {
						break
					}
				}
			}
			break
		}
	case 'a':
		{
			ch2, err := p.getChar()
			if err != nil {
				return 0, err
			}
			if p.encoding < FormatUtf16BE {
				if ch2 == 'm' {
					// checks that the sequence matcher &amp;
					if err := checkSeq("p;", offset, 1); err != nil {
						return 0, err
					}
					inc = 5
					val = '&'
				} else if ch2 == 'p' {
					// checks that the sequence matcher &apos;
					if err := checkSeq("os;", offset, 1); err != nil {
						return 0, err
					}
					inc = 5
					val = '\''
				}
			} else {
				if ch2 == 'm' {
					// checks that the sequence matcher &amp;
					if err := checkSeq("p;", offset, 2); err != nil {
						return 0, err
					}
					inc = 10
					val = '&'
				} else if ch2 == 'p' {
					// checks that the sequence matcher &apos;
					if err := checkSeq("os;", offset, 2); err != nil {
						return 0, err
					}
					inc = 12
					val = '\''
				}
			}
			break
		}
	case 'q':
		{
			if p.encoding < FormatUtf16BE {
				if err := checkSeq("uot;", offset, 1); err != nil {
					return 0, err
				}
				inc = 6
			} else {
				if err := checkSeq("uot;", offset, 2); err != nil {
					return 0, err
				}
				inc = 12
			}
			val = '"'
			break
		}
	case 'g', 'l':
		{
			if p.encoding < FormatUtf16BE {
				if err := checkSeq("t;", offset, 1); err != nil {
					return 0, err
				}
				inc = 4
			} else {
				if err := checkSeq("t;", offset, 2); err != nil {
					return 0, err
				}
				inc = 8
			}
			if ch == 'g' {
				val = '>'
			} else {
				val = '<'
			}
			break
		}
	}
	return int64(int(val) | (inc << 32)), nil
}

func (p *VtdParser) getCharUnit(offset int) (int32, error) {
	if p.encoding <= 2 {
		return int32(p.xmlDoc[offset] & 0xff), nil
	} else if p.encoding < FormatUtf16BE {
		ch, err := p.reader.Decode(int32(offset))
		if err != nil {
			return 0, err
		}
		return int32(ch), nil
	} else if p.encoding == FormatUtf16BE {
		return int32(p.xmlDoc[offset])<<8 | int32(p.xmlDoc[offset+1]), nil
	} else {
		return int32(p.xmlDoc[offset+1])<<8 | int32(p.xmlDoc[offset]), nil
	}
}

func (p *VtdParser) matchUrl(a, b Url) (bool, error) {
	for a.offset < a.offset+a.length &&
		b.offset < b.offset+b.length {
		chA, err := p.getCharResolved(a.offset)
		if err != nil {
			return false, err
		}
		chB, err := p.getCharResolved(b.offset)
		if err != nil {
			return false, err
		}
		if chA != chB {
			return false, err
		}
		a.offset += int(chA >> 32)
		b.offset += int(chB >> 32)
	}
	if a.offset == a.offset+a.length &&
		b.offset == b.offset+b.length {
		return true, nil
	}
	return false, nil
}
