package common

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	charSlice = []uint32{
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
		'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
		'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'/', '–', '?', ':', '(', ')', '.', ',', '‘', '+',
	}
)

func Test_XmlChar_IsValidChar_Success(t *testing.T) {
	x := NewXmlChar()
	for _, ch := range charSlice {
		assert.True(t, x.IsValidChar(ch), fmt.Sprintf("Invalid character %c", ch))
	}
}
