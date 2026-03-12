package gofig

import (
	"io"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
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

type YamlSource struct {
	Root    string
	records Records
}

func NewYamlSource(file string, root string) (*YamlSource, error) {
	data, err := os.ReadFile(file)

	if err != nil {
		return nil, err
	}

	records, err := parseYAML(data)
	if err != nil {
		return nil, err
	}

	return &YamlSource{
		Root:    root,
		records: records,
	}, nil
}

func parseYAML(data []byte) (Records, error) {
	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	records := make(Records)
	flatten(raw, "", records)
	return records, nil
}

func flatten(m map[string]any, prefix string, out Records) {
	for k, v := range m {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}

		switch val := v.(type) {
		case map[string]any:
			flatten(val, key, out)
		default:
			out[key] = val
		}
	}
}

func (this *YamlSource) Read(path string) any {
	return this.records[path]
}

var _ Source = (*YamlSource)(nil)
