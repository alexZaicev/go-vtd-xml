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

func (p *VtdParser) getNextProcessStateFromChar(ch uint32) (State, error) {
	if ch == '<' {
		if p.ws {
			if err := p.recordWhiteSpace(); err != nil {
				return StateInvalid, err
			}
		}
		return StateLtSeen, nil
	}
	if ch == '&' {
		ch, err := p.entityIdentifier()
		if err != nil {
			return StateInvalid, err
		}
		if !p.xmlChar.IsValidChar(ch) {
			return StateInvalid, erroring.NewParseError(erroring.InvalidChar, p.fmtLine(), nil)
		}
		return StateText, nil
	}
	if ch == ']' {
		// skip all ] chars
		for p.reader.SkipChar(']') {
		}
		if p.reader.SkipChar('>') {
			return StateInvalid, erroring.NewParseError("]]> sequence in text content", p.fmtLine(), nil)
		}
		return StateText, nil
	}
	if p.xmlChar.IsContentChar(ch) {
		return StateText, nil
	}
	return StateInvalid, erroring.NewParseError(erroring.InvalidChar, p.fmtLine(), nil)
}
