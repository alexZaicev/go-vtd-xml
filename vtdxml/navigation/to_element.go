package navigation

import (
	"github.com/alexZaicev/go-vtd-xml/vtdxml/common"
	"github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"
)

func (n *VtdNav) ToElement(dir Direction) (bool, error) {
	switch dir {
	case Root:
		return n.toRootElement()
	case Parent:
		return n.toParentElement()
	case FirstChild, LastChild:
		return n.toChildElement(dir)
	case NextSibling, PrevSibling:
		return n.toSiblingElement(dir)
	default:
		return false, erroring.NewNavigationError("illegal navigation options", nil)
	}
}

func (n *VtdNav) toRootElement() (bool, error) {
	if n.context[0] != 0 {
		n.context[0] = 0
	}
	n.atTerminal = false
	return true, nil
}

func (n *VtdNav) toParentElement() (bool, error) {
	if n.atTerminal == true {
		n.atTerminal = false
		return true, nil
	}
	if n.context[0] > 0 {
		n.context[n.context[0]] = -1
		n.context[0]--
		return true, nil
	} else if n.context[0] == 0 {
		n.context[0] = -1
		return true, nil
	} else {
		return false, nil
	}
}

func (n *VtdNav) toChildElement(dir Direction) (bool, error) {
	if n.atTerminal {
		return false, nil
	}
	switch n.context[0] {
	case -1:
		{
			n.context[0] = 0
			return true, nil
		}
	case 0:
		{
			if n.l1Buffer.GetSize() == 0 {
				return false, nil
			}
			n.context[0] = 1
			if dir == FirstChild {
				n.l1index = 0
			} else {
				n.l1index = n.l1Buffer.GetSize() - 1
			}
			val, err := n.l1Buffer.Upper32At(n.l1index)
			if err != nil {
				return false, err
			}
			n.context[1] = val
			return true, nil
		}
	case 1:
		{
			lower, err := n.l1Buffer.Lower32At(n.l1index)
			if err != nil {
				return false, err
			}
			n.l2lower = lower
			if n.l2lower == -1 {
				return false, nil
			}
			n.context[0] = 2
			n.l2upper = int32(n.l2Buffer.GetSize() - 1)
			for i := n.l1index + 1; i < n.l1Buffer.GetSize(); i++ {
				lower, err = n.l1Buffer.Lower32At(i)
				if err != nil {
					return false, err
				}
				if uint32(lower) != uint32(0xffffffff) {
					n.l2upper = lower - 1
					break
				}
			}
			if dir == FirstChild {
				n.l2index = int(n.l2lower)
			} else {
				n.l2index = int(n.l2upper)
			}
			upper, err := n.l2Buffer.Upper32At(n.l2index)
			if err != nil {
				return false, nil
			}
			n.context[2] = upper
			return true, nil
		}
	case 2:
		{
			lower, err := n.l2Buffer.Lower32At(n.l2index)
			if err != nil {
				return false, err
			}
			n.l3lower = lower
			if n.l3lower == -1 {
				return false, err
			}
			n.context[0] = 3
			n.l3upper = int32(n.l3Buffer.GetSize() - 1)
			for i := n.l2index + 1; i < n.l2Buffer.GetSize(); i++ {
				lower, err = n.l2Buffer.Lower32At(i)
				if err != nil {
					return false, err
				}
				if uint32(lower) != uint32(0xffffffff) {
					n.l3upper = lower - 1
					break
				}
			}
			if dir == FirstChild {
				n.l3index = int(n.l3lower)
			} else {
				n.l3index = int(n.l3upper)
			}
			upper, err := n.l2Buffer.Upper32At(n.l2index)
			if err != nil {
				return false, nil
			}
			n.context[3] = upper
			return true, nil
		}
	default:
		if dir == FirstChild {
			return n.toFirstChild()
		} else {
			return n.toLastChild()
		}
	}
}

func (n *VtdNav) toFirstChild() (bool, error) {
	index := int(n.context[n.context[0]] + 1)
	for index < n.vtdBuffer.GetSize() {
		tokenType, err := n.GetTokenType(index)
		if err != nil {
			return false, nil
		}
		depth, err := n.GetTokenDepth(index)
		if err != nil {
			return false, err
		}
		if common.Token(tokenType) == common.TokenStartingTag {
			if depth <= n.context[0] {
				return false, nil
			} else if depth == (n.context[0] + 1) {
				n.context[0] += 1
				n.context[n.context[0]] = int32(index)
				return true, nil
			}
		}
		index++
	}
	return false, nil
}

func (n *VtdNav) toLastChild() (bool, error) {
	index := int(n.context[n.context[0]] + 1)
	lastIndex := -1
	for index < n.vtdBuffer.GetSize() {
		tokenType, err := n.GetTokenType(index)
		if err != nil {
			return false, nil
		}
		depth, err := n.GetTokenDepth(index)
		if err != nil {
			return false, err
		}
		if common.Token(tokenType) == common.TokenStartingTag {
			if depth <= n.context[0] {
				break
			} else if depth == (n.context[0] + 1) {
				lastIndex = index
			}
		}
		index++
	}
	if lastIndex == -1 {
		return false, nil
	} else {
		n.context[0] += 1
		n.context[n.context[0]] = int32(lastIndex)
		return true, nil
	}
}

func (n *VtdNav) toSiblingElement(dir Direction) (bool, error) {
	if n.atTerminal {
		return false, nil
	}
	switch n.context[0] {
	case -1, 0:
		return false, nil
	case 1:
		{
			if dir == NextSibling {
				if n.l1index+1 >= n.l1Buffer.GetSize() {
					return false, nil
				}
				n.l1index++
			} else {
				if n.l1index-1 < 0 {
					return false, nil
				}
				n.l1index--
			}
			upper, err := n.l1Buffer.Upper32At(n.l1index)
			if err != nil {
				return false, err
			}
			n.context[1] = upper
			return true, nil
		}
	case 2:
		{
			if dir == NextSibling {
				if n.l2index+1 > int(n.l2upper) {
					return false, nil
				}
				n.l2index++
			} else {
				if n.l2index-1 < int(n.l2lower) {
					return false, nil
				}
				n.l2index--
			}
			upper, err := n.l1Buffer.Upper32At(n.l1index)
			if err != nil {
				return false, err
			}
			n.context[2] = upper
			return true, nil
		}
	case 3:
		{
			if dir == NextSibling {
				if n.l3index+1 > int(n.l3upper) {
					return false, nil
				}
				n.l3index++
			} else {
				if n.l3index-1 < int(n.l3lower) {
					return false, nil
				}
				n.l3index--
			}
			upper, err := n.l1Buffer.Upper32At(n.l1index)
			if err != nil {
				return false, err
			}
			n.context[3] = upper
			return true, nil
		}
	default:
		{
			if dir == NextSibling {
				return n.toNextSibling()
			} else {
				return n.toPrevSibling()
			}
		}
	}
}

func (n *VtdNav) toNextSibling() (bool, error) {
	index := int(n.context[n.context[0]] + 1)
	for index < n.vtdBuffer.GetSize() {
		tokenType, err := n.GetTokenType(index)
		if err != nil {
			return false, nil
		}
		depth, err := n.GetTokenDepth(index)
		if err != nil {
			return false, err
		}
		if common.Token(tokenType) == common.TokenStartingTag {
			if depth <= n.context[0] {
				return false, nil
			} else if depth == n.context[0] {
				n.context[n.context[0]] = int32(index)
				return true, nil
			}
		}
		index++
	}
	return false, nil
}

func (n *VtdNav) toPrevSibling() (bool, error) {
	index := int(n.context[n.context[0]] - 1)
	for index > int(n.context[n.context[0]-1]) {
		tokenType, err := n.GetTokenType(index)
		if err != nil {
			return false, nil
		}
		depth, err := n.GetTokenDepth(index)
		if err != nil {
			return false, err
		}
		if common.Token(tokenType) == common.TokenStartingTag {
			if depth == n.context[0] {
				n.context[n.context[0]] = int32(index)
				return true, nil
			}
		}
		index--
	}
	return false, nil
}
