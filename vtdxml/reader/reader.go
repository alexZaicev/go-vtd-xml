package reader

type Reader interface {
	GetChar() (int32, error)
	GetLongChar(offset int32) (int64, error)
	SkipChar(ch int32) (bool, error)
	Decode(offset int32) (string, error)
}
