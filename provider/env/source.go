package env

import (
	"os"
	"strings"

	"github.com/masnyjimmy/gofig"
)

func defaultTransform(path string) string {
	path = strings.ToUpper(path)
	path = strings.ReplaceAll(path, ".", "_")
	return path
}

type Source struct {
	p *Provider
}

func (s *Source) Read(path string) (any, error) {

	if s.p.prefix != "" {
		path = s.p.prefix + "." + path
	}

	if s.p.pathTransformFn == nil {
		path = defaultTransform(path)
	} else {
		path = s.p.pathTransformFn(path)
	}

	val, has := os.LookupEnv(path)

	if !has {
		return nil, nil
	}

	return val, nil
}

var _ gofig.Source = (*Source)(nil)
