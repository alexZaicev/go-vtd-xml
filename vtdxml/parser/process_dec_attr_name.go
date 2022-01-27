package parser

import "github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"

func (p *VtdParser) processDecAttrName() (State, error) {
	if p.currentChar == 'v' && p.reader.SkipCharSeq("ersion") {
		if err := p.nextCharAfterWs(); err != nil {
			return StateInvalid, err
		}
		if p.currentChar != '=' {
			return StateInvalid, erroring.NewParseError(erroring.InvalidChar, p.fmtLine(), nil)
		}
		if p.singleByteEncoding {
			if err := p.writeVtd(TokenDecAttrName, p.lastOffset-1, 7, p.depth); err != nil {
				return StateInvalid, err
			}
		} else {
			if err := p.writeVtd(TokenDecAttrName, (p.lastOffset-2)>>2, 7, p.depth); err != nil {
				return StateInvalid, err
			}
		}
	} else {
		return StateInvalid, erroring.NewParseError("declaration should be version", p.fmtLine(), nil)
	}

}
