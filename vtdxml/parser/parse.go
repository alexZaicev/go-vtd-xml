package parser

import (
	"errors"

	"github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"
)

// Parse function generates VTD tokens and location cache info.
// If namespace awareness is set to true, VTDGen conforms to XML
// namespace 1.0 spec
func (p *VtdParser) Parse() error {
	if err := p.decideEncoding(); err != nil {
		return err
	}
	if err := p.writeVtd(TokenDocument, 0, 0, p.depth); err != nil {
		return err
	}

	parserState := StateDocStart
	var err error

	for {
		switch parserState {
		case StateLtSeen:
			parserState, err = p.processLtSeen()
			break
		case StateStartTag:
			parserState, err = p.processStartTag()
			break
		case StateEndTag:
			parserState, err = p.processEndTag()
			break
		case StateAttrName:
			parserState, err = p.processAttrName()
			break
		case StateAttrVal:
			parserState, err = p.processAttrVal()
			break
		case StateText:
			parserState, err = p.processText()
			break
		case StateDocStart:
			parserState, err = p.processDocStart()
			break
		case StateDocEnd:
		case StatePiTag:
		case StatePiVal:
		case StateDecAttrName:
		case StateStartComment:
		case StateEndComment:
		case StateCdata:
		case StateDocType:
			parserState, err = p.processDocType()
			break
		case StateEndPi:
		default:
			return erroring.NewParseError(
				"invalid parser state", p.fmtLine(), nil,
			)
		}

		if errors.As(err, &erroring.EOFErrorType) && parserState == StateDocEnd {
			if err := p.finishUp(); err != nil {
				return erroring.NewInternalError("failed to finish-up document parsing", err)
			}
		} else if err != nil {
			return err
		}
	}
}

// finishUp function writes the remaining portion of LC info
func (p *VtdParser) finishUp() error {
	var err error
	if p.shallowDepth {
		if p.lastDepth == 1 {
			err = p.l1Buffer.Append(int64((p.lastL1Index << 32) | 0xFFFF))
		} else if p.lastDepth == 2 {
			err = p.l2Buffer.Append(int64((p.lastL2Index << 32) | 0xFFFF))
		}
	} else {
		if p.lastDepth == 1 {
			err = p.l1Buffer.Append(int64((p.lastL1Index << 32) | 0xFFFF))
		} else if p.lastDepth == 2 {
			err = p.l2Buffer.Append(int64((p.lastL2Index << 32) | 0xFFFF))
		} else if p.lastDepth == 3 {
			err = p.l3Buffer.Append(int64((p.lastL3Index << 32) | 0xFFFF))
		} else if p.lastDepth == 4 {
			err = p.l4Buffer.Append(int64((p.lastL4Index << 32) | 0xFFFF))
		}
	}
	return err
}
