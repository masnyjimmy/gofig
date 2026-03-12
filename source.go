package gofig

import (
	"os"
	"strings"
)

type Source interface {
	Read(path string) any
}

type EnvSourceConfig struct {
	Prefix string
}

type EnvTransformFunc = func(path string) string

type EnvSource struct {
	Prefix    string
	Transform EnvTransformFunc
}

func (this *EnvSource) pathToKey(path string) string {
	if this.Prefix != "" {
		path = this.Prefix + "." + path
	}

	if this.Transform != nil {
		path = this.Transform(path)
	} else {
		path = strings.ToUpper(path)
		path = strings.ReplaceAll(path, ".", "_")
	}
	return path
}

func (this *EnvSource) Read(path string) any {
	path = this.pathToKey(path)

	val, has := os.LookupEnv(path)

	if !has {
		return nil
	}

	return val
}

var _ Source = (*EnvSource)(nil)
