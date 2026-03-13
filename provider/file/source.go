package yaml

import (
	"github.com/masnyjimmy/gofig"
)

type source struct {
	records gofig.Records
}

func (this *source) Read(path string) (any, error) {
	return this.records[path], nil
}

var _ gofig.Source = (*source)(nil)
