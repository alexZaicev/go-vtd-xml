package parser

import (
	"errors"
	"fmt"

	"github.com/alexZaicev/go-vtd-xml/vtdxml/common"
	"github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"
)

// Parse function generates VTD tokens and location cache info.
// If namespace awareness is set to true, VTDGen conforms to XML
// namespace 1.0 spec
func (p *VtdParser) Parse() error {
	if err := p.decideEncoding(); err != nil {
		return err
	}
	if err := p.writeVtd(common.TokenDocument, 0, 0, p.depth); err != nil {
		return err
	}

	parserState := StateDocStart

	var ps State
	var err error
	for {
		fmt.Printf("Starting process state: %d\n", parserState)
		fmt.Printf("Offset: %d LastOffset %d CurrentChar %d Length1 %d\n", p.offset, p.lastOffset, p.currentChar,
			p.length1)

		switch parserState {
		case StateDocType:
			ps, err = p.processDocType()
		case StateDocStart:
			ps, err = p.processDocStart()
		case StateDocEnd:
			ps, err = p.processDocEnd()
		case StateLtSeen:
			ps, err = p.processLtSeen()
		case StateTagStart:
			ps, err = p.processStartTag()
		case StateTagEnd:
			ps, err = p.processEndTag()
		case StateAttrName:
			ps, err = p.processAttrName()
		case StateAttrVal:
			ps, err = p.processAttrVal()
		case StateDecAttrName:
			ps, err = p.processDecAttrName()
		case StateText:
			ps, err = p.processText()
		case StatePiTag:
			ps, err = p.processPiTag()
		case StatePiVal:
			ps, err = p.processPiVal()
		case StatePiEnd:
			ps, err = p.processPiEnd()
		case StateStartComment:
			ps, err = p.processStartComment()
		case StateEndComment:
			ps, err = p.processEndComment()
		case StateCdata:
			ps, err = p.processCdata()
		default:
			return erroring.NewParseError(
				"invalid parser state", p.fmtLine(), nil,
			)
		}

		if errors.As(err, &erroring.EOFErrorType) && parserState == StateDocEnd {
			if err := p.finishUp(); err != nil {
				return erroring.NewInternalError("failed to finish-up document parsing", err)
			}
			return nil
		} else if err != nil {
			return err
		}

		parserState = ps
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
