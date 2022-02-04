package navigation

import "github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"

func (n *VtdNav) ToElement(dir Direction) (bool, error) {
	switch dir {
	case Root:
		return n.toRootElement()
	case Parent:
		return n.toParentElement()
	case FirstChild, LastChild:
		return n.toChildElement()
	case NextSibling, PrevSibling:
		return n.toSiblingElement()
	default:
		return false, erroring.NewNavigationError("illegal navigation options", nil)
	}
}

func (n *VtdNav) toRootElement() (bool, error) {
	return false, nil
}

func (n *VtdNav) toParentElement() (bool, error) {
	return false, nil
}

func (n *VtdNav) toChildElement() (bool, error) {
	return false, nil
}

func (n *VtdNav) toSiblingElement() (bool, error) {
	return false, nil
}
