package navigation

import (
	"github.com/alexZaicev/go-vtd-xml/vtdxml/common"
	"github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"
)

func (n *VtdNav) getCharResolved(offset int) (uint32, error) {
	ch, err := n.getChar(offset)
	if err != nil {
		return 0, err
	}
	if ch != '&' {
		return uint32(ch), nil
	}
	offset++
	ch2, err := n.getCharUnit(offset)
	if err != nil {
		return 0, err
	}
	offset++
	inc, val, err := n.resolveEntity(offset, ch2)
	if err != nil {
		return 0, err
	}
	return val | uint32(inc<<32), nil
}

func (n *VtdNav) resolveEntity(offset int, ch uint32) (uint64, uint32, error) {
	checkSeq := func(os int, seq string) error {
		for i, seqCh := range seq {
			ch, err := n.getCharUnit(os + i)
			if err != nil {
				return err
			}
			if int32(ch) != seqCh {
				return erroring.NewEntityError(erroring.IllegalBuiltInEntity)
			}
		}
		return nil
	}
	var inc uint64
	switch ch {
	case '#':
		{
			var value uint32
			ch2, err := n.getCharUnit(offset)
			if err != nil {
				return 0, 0, err
			}
			inc = 2
			if ch2 == 'x' {
				for {
					offset++
					inc++
					ch2, err = n.getCharUnit(offset)
					if err != nil {
						return 0, 0, err
					}
					if ch2 >= '0' && ch2 <= '9' {
						value = (value << 4) + (ch2 - '0')
					} else if ch2 >= 'a' && ch2 <= 'f' {
						value = (value << 4) + (ch2 - 'a' + 10)
					} else if ch2 >= 'A' && ch2 <= 'F' {
						value = (value << 4) + (ch2 - 'A' + 10)
					} else if ch2 == ';' {
						break
					} else {
						return 0, 0, erroring.NewEntityError("illegal char following &#x")
					}
				}
			} else {
				for {
					if ch2 >= '0' && ch2 <= '9' {
						value = value*10 + (ch2 - '0')
					} else if ch2 == ';' {
						break
					} else {
						return 0, 0, erroring.NewEntityError("illegal char following &#x")
					}
					ch2, err = n.getCharUnit(offset)
					if err != nil {
						return 0, 0, err
					}
				}
			}
			if !n.xmlChar.IsValidChar(value) {
				return 0, 0, erroring.NewEntityError(erroring.InvalidChar)
			}
			return inc, value, nil
		}
	case 'a':
		{
			// checks that the sequence matcher &amp;
			ch2, err := n.getCharUnit(offset)
			if err != nil {
				return 0, 0, err
			}
			if ch2 == 'm' {
				// checks that the sequence matcher &amp;
				if err := checkSeq(offset, "p;"); err != nil {
					return 0, 0, err
				}
				return 5, '&', nil
			} else if ch2 == 'p' {
				// checks that the sequence matcher &apos;
				if err := checkSeq(offset, "os;"); err != nil {
					return 0, 0, err
				}
				return 6, '\'', nil
			} else {
				return 0, 0, erroring.NewEntityError(erroring.IllegalBuiltInEntity)
			}
		}
	case 'q':
		{
			// checks that the sequence matcher &quot;
			if err := checkSeq(offset, "uot;"); err != nil {
				return 0, 0, err
			}
			return 6, '"', nil
		}
	case 'l', 'g':
		{
			// checks that the sequence matcher &gt; or &lt;
			if err := checkSeq(offset, "t;"); err != nil {
				return 0, 0, err
			}
			val := '>'
			if ch == 'l' {
				val = '<'
			}
			return 4, uint32(val), nil
		}
	default:
		return 0, 0, erroring.NewEntityError("illegal entity character")
	}
}

func (n *VtdNav) getCharUnit(offset int) (uint32, error) {
	if n.encoding == common.FormatAscii ||
		n.encoding == common.FormatIso88591 ||
		n.encoding == common.FormatUtf8 {
		b, err := n.xmlBuffer.ByteAt(offset)
		if err != nil {
			return 0, err
		}
		return uint32(b & 0xFF), nil
	}

	if (n.encoding > common.FormatIso88591 && n.encoding < common.FormatIso885916) ||
		(n.encoding >= common.FormatWin1250 && n.encoding <= common.FormatWin1258) {
		return n.decode(offset)
	}

	if n.encoding == common.FormatUtf16BE {
		b1, err := n.xmlBuffer.ByteAt(offset << 1)
		if err != nil {
			return 0, err
		}
		b2, err := n.xmlBuffer.ByteAt((offset << 1) + 1)
		if err != nil {
			return 0, err
		}
		return uint32(b1)<<8 | uint32(b2), nil
	} else {
		b1, err := n.xmlBuffer.ByteAt(offset << 1)
		if err != nil {
			return 0, err
		}
		b2, err := n.xmlBuffer.ByteAt((offset << 1) + 1)
		if err != nil {
			return 0, err
		}
		return uint32(b2)<<8 | uint32(b1), nil
	}
}

func (n *VtdNav) getChar(offset int) (uint64, error) {
	b, err := n.xmlBuffer.ByteAt(offset)
	if err != nil {
		return 0, err
	}
	switch n.encoding {
	case common.FormatAscii:
		{
			if b == '\r' {
				b2, ifErr := n.xmlBuffer.ByteAt(offset + 1)
				if ifErr != nil {
					return 0, nil
				}
				if b2 == '\n' {
					return uint64(int(b2) | (2 << 32)), nil
				} else {
					return uint64(int('\n') | (1 << 32)), nil
				}
			}

			return uint64(int(b) | (1 << 32)), nil
		}
	case common.FormatIso88591:
		{
			if b == '\r' {
				b2, ifErr := n.xmlBuffer.ByteAt(offset + 1)
				if ifErr != nil {
					return 0, nil
				}
				if b2 == '\n' {
					return uint64(int(b2) | (2 << 32)), nil
				} else {
					return uint64(int('\n') | (1 << 32)), nil
				}
			}

			return uint64((int(b) & 0xFF) | (1 << 32)), nil
		}
	case common.FormatUtf8:
		{
			if b == '\r' {
				b2, ifErr := n.xmlBuffer.ByteAt(offset + 1)
				if ifErr != nil {
					return 0, nil
				}
				if b2 == '\n' {
					return uint64(int(b2) | (2 << 32)), nil
				} else {
					return uint64(int('\n') | (1 << 32)), nil
				}
			}

			return uint64(int(b) | (1 << 32)), nil
		}
	case common.FormatUtf16BE:
		return n.getCharUtf16BE(offset)
	case common.FormatUtf16LE:
		return n.getCharUtf16LE(offset)
	default:
		return n.getChar4OtherEncodings(offset)
	}
}

func (n *VtdNav) getChar4OtherEncodings(offset int) (uint64, error) {
	if n.encoding > common.FormatWin1258 {
		return 0, erroring.NewEncodingError("unknown encoding")
	}
	// TODO implement
	return 0, nil
}

func (n *VtdNav) getCharUtf16BE(offset int) (uint64, error) {
	// TODO implement
	return 0, nil
}

func (n *VtdNav) getCharUtf16LE(offset int) (uint64, error) {
	// TODO implement
	return 0, nil
}
