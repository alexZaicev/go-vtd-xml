package common

import "github.com/alexZaicev/go-vtd-xml/vtdxml/buffer"

const (
	DefaultWidth     = 0
	DefaultHashWidth = 1
	DefaultMask      = 1
)

type Option func(*IntHash)

type IntHash struct {
	storage                    []buffer.FastIntBuffer
	maxDepth, hashWidth, width int
	mask                       int
}

func WithSize(size int) Option {
	return func(hash *IntHash) {
		width := determineHashWidth(size)

		hash.width = width
		hash.hashWidth = 1 << width
		hash.storage = make([]buffer.FastIntBuffer, 0, hash.hashWidth)
	}
}

func NewIntHash(opts ...Option) *IntHash {
	hash := &IntHash{
		width:     DefaultWidth,
		hashWidth: DefaultHashWidth,
		mask:      DefaultMask,
		storage:   make([]buffer.FastIntBuffer, 0, DefaultHashWidth),
	}

	for _, opt := range opts {
		opt(hash)
	}
	return hash
}

func determineHashWidth(i int) int {
	if i < (1 << 8) {
		return 3
	}
	if i < (1 << 9) {
		return 4
	}
	if i < (1 << 10) {
		return 5
	}
	if i < (1 << 11) {
		return 6
	}
	if i < (1 << 12) {
		return 7
	}
	if i < (1 << 13) {
		return 8
	}
	if i < (1 << 14) {
		return 9
	}
	if i < (1 << 15) {
		return 10
	}
	if i < (1 << 16) {
		return 11
	}
	if i < (1 << 17) {
		return 12
	}
	if i < (1 << 18) {
		return 13
	}
	if i < (1 << 19) {
		return 14
	}
	if i < (1 << 20) {
		return 15
	}
	if i < (1 << 21) {
		return 16
	}
	if i < (1 << 22) {
		return 17
	}
	if i < (1 << 23) {
		return 18
	}
	if i < (1 << 25) {
		return 19
	}
	if i < (1 << 27) {
		return 20
	}
	if i < (1 << 29) {
		return 21
	}
	return 22
}
