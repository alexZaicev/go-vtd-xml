package reader

type Reader interface {
	GetChar() (uint32, error)
	GetLongChar(offset int32) (uint64, error)
	SkipChar(ch uint32) bool
	SkipCharSeq(seq string) bool
	Decode(offset int32) (uint32, error)
}
