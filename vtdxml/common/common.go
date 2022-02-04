package common

type Token int64

const (
	TokenStartingTag Token = iota
	TokenEndingTag
	TokenAttrName
	TokenAttrNs
	TokenAttrVal
	TokenCharacterData
	TokenComment
	TokenPiName
	TokenPiVal
	TokenDecAttrName
	TokenDecAttrVal
	TokenCdataVal
	TokenDtdVal
	TokenDocument
)

const (
	MaskTokenType       uint64 = 0xF000000000000000
	MaskTokenOffset     uint64 = 0x000000003FFFFFFF
	MaskTokenQnLength   uint64 = 0x000007FF00000000
	MaskTokenFullLength uint64 = 0x000FFFFF00000000
	MaskTokenPreLength  uint64 = 0x000FF80000000000
	MaskTokenDepth      uint64 = 0x0FF0000000000000
)

type FormatEncoding int

const (
	FormatAscii FormatEncoding = iota
	FormatIso88591
	FormatIso88592
	FormatIso88593
	FormatIso88594
	FormatIso88595
	FormatIso88596
	FormatIso88597
	FormatIso88598
	FormatIso88599
	FormatIso885910
	FormatIso885911
	FormatIso885912
	FormatIso885913
	FormatIso885914
	FormatIso885915
	FormatIso885916
	FormatUtf16BE
	FormatUtf16LE
	FormatUtf8
	FormatWin1250
	FormatWin1251
	FormatWin1252
	FormatWin1253
	FormatWin1254
	FormatWin1255
	FormatWin1256
	FormatWin1257
	FormatWin1258
)
