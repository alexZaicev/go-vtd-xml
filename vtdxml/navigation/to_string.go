package navigation

import (
	"bytes"

	"github.com/alexZaicev/go-vtd-xml/vtdxml/common"
)

func (n *VtdNav) ToStringAtIndex(index int) (string, error) {
	tokenType, err := n.GetTokenType(index)
	if err != nil {
		return "", err
	}
	if common.Token(tokenType) != common.TokenCharacterData &&
		common.Token(tokenType) != common.TokenAttrVal {
		return n.ToRawStringAtIndex(index)
	}
	length, err := n.GetTokenLength(index)
	if err != nil {
		return "", err
	}
	offset, err := n.GetTokenOffset(index)
	if err != nil {
		return "", err
	}
	return n.ToStringAtRange(offset, length)
}

func (n *VtdNav) ToStringAtRange(offset, length int32) (string, error) {
	var buffer bytes.Buffer

	endOffset := offset + length
	for i := offset; i < endOffset; i++ {
		ch, err := n.getCharResolved(int(i))
		if err != nil {
			return "", err
		}
		buffer.WriteByte(byte(ch))
	}

	return buffer.String(), nil
}

func (n *VtdNav) ToRawStringAtIndex(index int) (string, error) {
	tokenType, err := n.GetTokenType(index)
	if err != nil {
		return "", err
	}
	var length int32
	if common.Token(tokenType) == common.TokenStartingTag ||
		common.Token(tokenType) == common.TokenAttrName ||
		common.Token(tokenType) == common.TokenAttrNs {
		l, ifErr := n.GetTokenLength(index)
		if ifErr != nil {
			return "", ifErr
		}
		length = l & 0xFFFF
	} else {
		l, ifErr := n.GetTokenLength(index)
		if ifErr != nil {
			return "", ifErr
		}
		length = l
	}
	offset, err := n.GetTokenOffset(index)
	if err != nil {
		return "", err
	}
	return n.ToRawStringAtRange(offset, length)
}

func (n *VtdNav) ToRawStringAtRange(offset, length int32) (string, error) {
	var buffer bytes.Buffer

	endOffset := offset + length
	for i := offset; i < endOffset; i++ {
		ch, err := n.getChar(int(i))
		if err != nil {
			return "", err
		}
		buffer.WriteByte(byte(ch))
	}

	return buffer.String(), nil
}
