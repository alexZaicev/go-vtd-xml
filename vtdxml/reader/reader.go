package reader

type Reader interface {
	GetChar() (uint32, error)
	GetLongCharAt(offset int32) (uint64, error)
	SkipChar(ch uint32) bool
	SkipCharSeq(seq string) bool
	GetCharAt(offset int32) (uint32, error)
	GetOffset() int
	SetOffset(offset int)
}
