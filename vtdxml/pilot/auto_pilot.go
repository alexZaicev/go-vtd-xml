package pilot

import (
	"github.com/alexZaicev/go-vtd-xml/vtdxml/buffer"
	"github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"
	"github.com/alexZaicev/go-vtd-xml/vtdxml/expresion"
	"github.com/alexZaicev/go-vtd-xml/vtdxml/navigation"
)

type Option func(*AutoPilot) error

type AutoPilot struct {
	nav                        navigation.Nav
	name                       string
	iterType                   IterationType
	ft, special, enableCaching bool
	size                       int
	intBuffer                  buffer.IntBuffer
	exp                        expresion.Expression
	symbolMap                  map[string]expresion.Expression
}

func WithNavigation(nav navigation.Nav) Option {
	return func(p *AutoPilot) error {
		if nav == nil {
			return erroring.NewInvalidArgumentError("nav", erroring.CannotBeNil, nil)
		}
		p.nav = nav
		return nil
	}
}

func WithCaching(enable bool) Option {
	return func(p *AutoPilot) error {
		p.enableCaching = enable
		return nil
	}
}

func NewAutoPilot(opts ...Option) (*AutoPilot, error) {
	m := &AutoPilot{
		iterType:      Undefined,
		ft:            true,
		enableCaching: true,
	}

	for _, opt := range opts {
		if optErr := opt(m); optErr != nil {
			return nil, optErr
		}
	}
	return m, nil
}
