package navigation

import (
	"github.com/alexZaicev/go-vtd-xml/vtdxml/common"
	"github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"
)

type Direction int

const (
	Root Direction = iota
	Parent
	FirstChild
	LastChild
	NextSibling
	PrevSibling
)

func (n *VtdNav) GetCurrentIndex() (int32, error) {
	if n.atTerminal {
		return n.ln, nil
	}
	switch n.context[0] {
	case -1:
		return 0, nil
	case 0:
		return n.rootIndex, nil
	}
	idx := int(n.context[0])
	if idx < 0 || idx >= len(n.context) {
		return 0, erroring.NewInternalError(erroring.IndexOutOfRange, nil)
	}
	return n.context[idx], nil
}

func (n *VtdNav) GetTokenType(index int) (int32, error) {
	val, err := n.vtdBuffer.LongAt(index)
	if err != nil {
		return 0, err
	}
	res := (uint64(val) & common.MaskTokenType) >> 60
	return int32(res) & 0xF, nil
}

func (n *VtdNav) GetTokenOffset(index int) (int32, error) {
	val, err := n.vtdBuffer.LongAt(index)
	if err != nil {
		return 0, err
	}
	return int32(uint64(val) & common.MaskTokenOffset), nil
}

func (n *VtdNav) GetTokenDepth(index int) (int32, error) {
	val, err := n.vtdBuffer.LongAt(index)
	if err != nil {
		return 0, err
	}
	res := (uint64(val) & common.MaskTokenDepth) >> 52
	if res != 255 {
		return int32(res), nil
	}
	return -1, nil
}

func (n *VtdNav) GetTokenLength(index int) (int32, error) {
	tokenType, err := n.GetTokenType(index)
	if err != nil {
		return 0, err
	}
	switch common.Token(tokenType) {
	case common.TokenAttrName, common.TokenAttrNs, common.TokenStartingTag:
		{
			val, caseErr := n.vtdBuffer.LongAt(index)
			if caseErr != nil {
				return 0, caseErr
			}
			var length uint64
			if n.nsAware {
				a := (uint64(val) & common.MaskTokenQnLength) >> 32
				b := ((uint64(val) & common.MaskTokenPreLength) >> 32) << 5
				length = a | b
			} else {
				length = (uint64(val) & common.MaskTokenQnLength) >> 32
			}
			return int32(length), nil
		}
	case common.TokenCharacterData, common.TokenCdataVal, common.TokenComment:
		{
			var length uint64
			depth, caseErr := n.GetTokenDepth(index)
			if caseErr != nil {
				return 0, caseErr
			}
			var offset int32
			var tokenTypeNew, tokenOffsetNew, tokenDepthNew int32
			tokenDepthNew = depth
			tokenTypeNew = tokenType

			for i := index; i < n.vtdBuffer.GetSize() &&
				depth == tokenDepthNew &&
				tokenType == tokenTypeNew &&
				offset == tokenOffsetNew; i++ {

				val, forErr := n.vtdBuffer.LongAt(i)
				if forErr != nil {
					return 0, forErr
				}
				length += (uint64(val) & common.MaskTokenFullLength) >> 32

				tokenOffsetNew, forErr = n.GetTokenOffset(i)
				if forErr != nil {
					return 0, forErr
				}
				tokenOffsetNew = int32((uint64(val) & common.MaskTokenFullLength) >> 32)
			}
			return int32(length), nil
		}
	}
	val, err := n.vtdBuffer.LongAt(index)
	if err != nil {
		return 0, err
	}
	length := (uint64(val) & common.MaskTokenFullLength) >> 32
	return int32(length), nil
}

// func (n *VtdNav) GetRawTokenLength(index int) (int, error) {
// 	val, err := n.vtdBuffer.LongAt(index)
// 	if err != nil {
// 		return 0, err
// 	}
// 	res := (uint64(val) & common.MaskTokenFullLength) >> 32
// 	return int(res), nil
// }

func (n *VtdNav) GetVtdBufferSize() int {
	return n.vtdBuffer.GetSize()
}

func (n *VtdNav) GetEncoding() common.FormatEncoding {
	return n.encoding
}
