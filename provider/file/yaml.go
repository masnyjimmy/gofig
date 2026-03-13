package file

import (
	"github.com/goccy/go-yaml"
	"github.com/masnyjimmy/gofig"
)

func parseYAML(data []byte) (gofig.Records, error) {
	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	records := make(gofig.Records)
	flatten(raw, "", records)
	return records, nil
}

func flatten(m map[string]any, prefix string, out gofig.Records) {
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
