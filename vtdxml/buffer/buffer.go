package buffer

// Buffer interface is the base interface with basic methods
// to meet required buffer functionality
type Buffer interface {
	GetSize() int
	SetSize(size int)
	Clear()
}
