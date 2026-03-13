package file

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/masnyjimmy/gofig"
)

type ParserFn = func(data []byte) (gofig.Records, error)

func getDefaultParser(ext string) (ParserFn, error) {
	switch ext {
	case ".yaml":
		return parseYAML, nil
	default:
		return nil, fmt.Errorf("Unsupported ext: %v", ext)
	}
}

type Provider struct {
	root     string
	files    []string
	parserFn ParserFn
}

func New() *Provider {
	return new(Provider)
}

func (p *Provider) SetParser(parser ParserFn) *Provider {
	p.parserFn = parser
	return p
}

func (p *Provider) SetRoot(root string) *Provider {
	p.root = root
	return p
}

func (p *Provider) SetFiles(filename ...string) *Provider {
	p.files = filename
	return p
}

func (p *Provider) AddFiles(filename ...string) *Provider {
	p.files = append(p.files, filename...)
	return p
}

func (p *Provider) Source() (gofig.Source, error) {
	for _, file := range p.files {
		data, err := os.ReadFile(file)

		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return nil, err
		}

		parse := p.parserFn

		if parse == nil {
			parse, err = getDefaultParser(filepath.Ext(file))

			if err != nil {
				return nil, fmt.Errorf("Unable to parse data: %w", err)
			}
		}

		records, err := parse(data)

		if err != nil {
			return nil, fmt.Errorf("Unable to parse data: %w", err)
		}

		return &source{
			records: records,
		}, nil
	}

	return nil, nil
}

var _ gofig.Provider = (*Provider)(nil)
