package parser

import (
	"fmt"

	"github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"
)

// validateSeq function to validate string sequence directly from the reader.Reader
func (p *VtdParser) validateSeq(seq string) error {
	for _, seqChar := range seq {
		c, err := p.reader.GetChar()
		if err != nil {
			return err
		}
		if c != uint32(seqChar) {
			return erroring.NewParseError(fmt.Sprintf("invalid char sequence in %s", seq),
				p.fmtLine(), nil)
		}
	}
	if p.depth < 0 {
		return erroring.NewParseError(fmt.Sprintf("wrong place for %s", seq),
			p.fmtLine(), nil)
	}
	return nil
}
