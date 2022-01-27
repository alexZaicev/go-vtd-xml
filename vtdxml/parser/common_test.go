package parser

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	version := "version"
	fmt.Println(uint32(version[0]))
	fmt.Println(version[1:])
	fmt.Println(version[1:3])

	utf8 := "utf-8"
	fmt.Println(utf8[1:4])

	utf16 := "utf-16"
	fmt.Println(utf16[4:6])
}
