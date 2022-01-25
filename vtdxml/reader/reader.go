package reader

type Reader interface {
	GetChar() (uint32, error)
	GetLongChar(offset uint32) (uint64, error)
	SkipChar(ch uint32) bool
	Decode(offset uint32) (string, error)
}
