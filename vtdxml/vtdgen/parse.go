package vtdgen

import (
	"io"

	"github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"
)

// Parse function generates VTD tokens and location cache info.
// If namespace awareness is set to true, VTDGen conforms to XML
// namespace 1.0 spec
func (g *VtdGen) Parse() error {
	if err := g.decideEncoding(); err != nil {
		return err
	}
	if err := g.writeVtd(TokenDocument, 0, 0, g.depth); err != nil {
		return err
	}

	parserState := StateDocStart
	var err error

	for {
		switch parserState {
		case StateLtSeen:
			parserState, err = g.processLtSeen()
			break
		case StateStartTag:
			parserState, err = g.processStartTag()
			break
		case StateEndTag:
		case StateAttrName:
		case StateAttrVal:
		case StateText:
		case StateDocStart:
		case StateDocEnd:
		case StatePiTag:
		case StatePiVal:
		case StateDecAttrName:
		case StateStartComment:
		case StateEndComment:
		case StateCdata:
		case StateDocType:
			parserState, err = g.processDocType()
			break
		case StateEndPi:
		default:
			return erroring.NewParseError(
				"invalid parser state", g.formatLineNumber(), nil,
			)
		}

		if err == io.EOF && parserState == StateDocEnd {
			if err := g.finishUp(); err != nil {
				return erroring.NewInternalError("failed to finish-up document parsing", err)
			}
		} else if err != nil {
			return err
		}
	}
}

// finishUp function writes the remaining portion of LC info
func (g *VtdGen) finishUp() error {
	// TODO implement
	return nil
}
