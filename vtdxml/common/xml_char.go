package common

type XmlChar struct {
	validCharSlice     []uint32
	spaceCharSlice     []uint32
	nameCharSlice      []uint32
	nameStartCharSlice []uint32
	contentCharSlice   []uint32
	specialCharSlice   []uint32
}

func NewXmlChar() *XmlChar {
	x := &XmlChar{
		validCharSlice:     make([]uint32, 0),
		spaceCharSlice:     make([]uint32, 0),
		nameCharSlice:      make([]uint32, 0),
		nameStartCharSlice: make([]uint32, 0),
		contentCharSlice:   make([]uint32, 0),
		specialCharSlice:   make([]uint32, 0),
	}

	// Set all valid characters
	rangeSlice := [][]uint32{{0x1, 0xD7FF}, {0xE000, 0xFFFD}, {0x10000, 0x10FFFF}}
	for _, i := range rangeSlice {
		for j := i[0]; j <= i[1]; j++ {
			// exclude surrogate blocks, 0xFFFE and 0xFFFF
			if (j >= 0xD800 && j <= 0xDFFF) || j == 0xFFFE || j == 0xFFFF {
				continue
			}
			x.validCharSlice = append(x.validCharSlice, j)
		}
	}

	// Set space characters
	for _, ch := range []uint32{0x20, 0x9, 0xD, 0xA} {
		x.spaceCharSlice = append(x.spaceCharSlice, ch)
	}

	// Set name start characters
	rangeSlice = [][]uint32{
		{':', ':'}, {'A', 'Z'}, {'_', '_'}, {'a', 'z'}, {0xC0, 0xD6}, {0xD8, 0xF6}, {0xF8, 0x2FF},
		{0x370, 0x37D}, {0x37F, 0x1FFF}, {0x200C, 0x200D}, {0x2070, 0x218F}, {0x2C00, 0x2FEF},
		{0x3001, 0xD7FF}, {0xF900, 0xFDCF}, {0xFDF0, 0xFFFD}, {0x10000, 0xEFFFF},
	}
	for _, i := range rangeSlice {
		for j := i[0]; j <= i[1]; j++ {
			x.nameStartCharSlice = append(x.nameStartCharSlice, j)
		}
	}

	// Set name characters
	rangeSlice = append(rangeSlice, []uint32{'-', '-'})
	rangeSlice = append(rangeSlice, []uint32{'.', '.'})
	rangeSlice = append(rangeSlice, []uint32{'0', '9'})
	rangeSlice = append(rangeSlice, []uint32{0xB7, 0xB7})
	rangeSlice = append(rangeSlice, []uint32{0x0300, 0x036F})
	rangeSlice = append(rangeSlice, []uint32{0x203F, 0x2040})
	for _, i := range rangeSlice {
		for j := i[0]; j <= i[1]; j++ {
			x.nameCharSlice = append(x.nameCharSlice, j)
		}
	}

	for _, i := range []uint32{'<', '&', ']'} {
		x.specialCharSlice = append(x.specialCharSlice, i)
	}

	return x
}

func (x *XmlChar) IsNameStartChar(ch uint32) bool {
	return x.exists(x.nameStartCharSlice, ch)
}

func (x *XmlChar) IsSpaceChar(ch uint32) bool {
	return x.exists(x.spaceCharSlice, ch)
}

func (x *XmlChar) IsValidChar(ch uint32) bool {
	return x.exists(x.validCharSlice, ch)
}

func (x *XmlChar) IsNameChar(ch uint32) bool {
	return x.exists(x.nameCharSlice, ch)
}

func (x *XmlChar) IsContentChar(ch uint32) bool {
	return !x.exists(x.specialCharSlice, ch) && x.exists(x.validCharSlice, ch)
}

func (x *XmlChar) exists(arr []uint32, ch uint32) bool {
	for _, i := range arr {
		if i == ch {
			return true
		}
	}
	return false
}
