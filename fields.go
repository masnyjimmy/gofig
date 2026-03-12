package gofig

import (
	"reflect"
	"strings"
	"time"
)

type FieldInfo struct {
	Required bool
	Default  any
}

func fieldInfo(f reflect.StructField) (path string, info FieldInfo, has bool) {
	tag, has := f.Tag.Lookup(LABEL)

	if !has {
		return "", FieldInfo{}, false
	}

	values := strings.Split(tag, ",")

	path = values[0]

	if len(values) > 1 {
		info.Required = false
		if def := values[1]; def != "-" {
			info.Default = def
		}
	} else {
		info.Required = true
	}

	return path, info, true
}

type Fields = []string

func parseStructToFields(t reflect.Type, current Fields, path string) Fields {

	for idx := range t.NumField() {
		f := t.Field(idx)
		ft := f.Type

		if f.Anonymous {
			continue
		}

		tag, has := f.Tag.Lookup(LABEL)

		if !has {
			continue
		}

		name, _, _ := strings.Cut(tag, ",")

		var currentPath string

		if path == "" {
			currentPath = name
		} else {
			currentPath = path + "." + name
		}

		if ft.Kind() == reflect.Struct && ft != reflect.TypeFor[time.Time]() && ft != reflect.TypeFor[time.Duration]() {
			current = parseStructToFields(f.Type, current, currentPath)
		} else {
			current = append(current, currentPath)
		}
	}

	return current
}

func GenerateFields[T any]() Fields {
	return parseStructToFields(reflect.TypeFor[T](), make(Fields, 0), "")
}
