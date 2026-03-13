package env

import "github.com/masnyjimmy/gofig"

type PathTranformFn = func(path string) string

type Provider struct {
	pathTransformFn PathTranformFn
	prefix          string
}

func (p *Provider) SetPathTransformer(fn PathTranformFn) *Provider {
	p.pathTransformFn = fn
	return p
}

func (p *Provider) SetPrefix(prefix string) *Provider {
	p.prefix = prefix
	return p
}

func (p *Provider) Source() (gofig.Source, error) {
	return &Source{
		p: p,
	}, nil
}

var _ gofig.Provider = (*Provider)(nil)
