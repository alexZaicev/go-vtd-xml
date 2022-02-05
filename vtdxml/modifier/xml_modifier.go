package modifier

import (
	"github.com/alexZaicev/go-vtd-xml/vtdxml/buffer"
	"github.com/alexZaicev/go-vtd-xml/vtdxml/common"
	"github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"
	"github.com/alexZaicev/go-vtd-xml/vtdxml/navigation"
)

type Option func(*XmlModifier) error

type XmlModifier struct {
	nav                    navigation.Nav
	longBuffer             buffer.LongBuffer
	objBuffer              buffer.ObjectBuffer
	insertHash, deleteHash *common.IntHash
	encoding               common.FormatEncoding
	charset                string
}

func WithNavigation(nav navigation.Nav) Option {
	return func(m *XmlModifier) error {
		if nav == nil {
			return erroring.NewInvalidArgumentError("nav", erroring.CannotBeNil, nil)
		}
		m.nav = nav
		m.insertHash = common.NewIntHash(common.WithSize(nav.GetVtdBufferSize()))
		m.deleteHash = common.NewIntHash(common.WithSize(nav.GetVtdBufferSize()))
		m.encoding = nav.GetEncoding()
		cs, err := getCharsetFromEncoding(m.encoding)
		if err != nil {
			return err
		}
		m.charset = cs
		return nil
	}
}

func NewXmlModifier(opts ...Option) (*XmlModifier, error) {
	m := &XmlModifier{}

	longBuffer, err := buffer.NewFastLongBuffer()
	if err != nil {
		return nil, err
	}
	m.longBuffer = longBuffer

	objBuffer, err := buffer.NewFastObjectBuffer()
	if err != nil {
		return nil, err
	}
	m.objBuffer = objBuffer

	for _, opt := range opts {
		if optErr := opt(m); optErr != nil {
			return nil, optErr
		}
	}
	return m, nil
}

func getCharsetFromEncoding(enc common.FormatEncoding) (string, error) {
	switch enc {
	case common.FormatAscii:
		return "ASCII", nil
	case common.FormatUtf8:
		return "UTF8", nil
	default:
		return "", erroring.NewModifyError("Master document encoding not yet supported by XML modifier")
	}
}
