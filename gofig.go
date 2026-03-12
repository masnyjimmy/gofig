package gofig

import (
	"fmt"
	"reflect"
	"time"
)

const LABEL = "conf"

type Records = map[string]any

type Gofig struct {
	Fields      Fields
	Records     Records
	TimeFormats []string
}

func New(fields Fields) *Gofig {
	return &Gofig{
		Fields:  fields,
		Records: make(Records),
	}
}

func (this *Gofig) SetTimeFormats(formats ...string) {
	this.TimeFormats = formats
}

func (this *Gofig) AddTimeFormats(formats ...string) {
	this.TimeFormats = append(this.TimeFormats, formats...)
}

func (this *Gofig) Read(src Source) {
	for _, path := range this.Fields {
		value := src.Read(path)

		if value != nil {
			this.Records[path] = value
		}
	}
}

func (this *Gofig) applyRecords(dest reflect.Value, path string) error {

	t := dest.Type()

	for idx := range t.NumField() {

		f := t.Field(idx)
		ft := f.Type

		name, fInfo, has := fieldInfo(f)

		if !has {
			continue
		}

		var currentPath string

		if path == "" {
			currentPath = name
		} else {
			currentPath = path + "." + name
		}

		if ft.Kind() == reflect.Struct && ft != reflect.TypeFor[time.Time]() {
			if err := this.applyRecords(dest.Field(idx), currentPath); err != nil {
				return err
			}
		} else {
			value, has := this.Records[currentPath]

			if has {
				if err := this.applyValue(dest.Field(idx), value); err != nil {
					return fmt.Errorf("Unable to apply value: %w", err)
				}
			} else if fInfo.Required {
				return fmt.Errorf("Missing required value: %v", name)
			} else {
				if err := this.applyValue(dest.Field(idx), fInfo.Default); err != nil {
					return fmt.Errorf("Unable to apply value: %w", err)
				}
			}
		}
	}

	return nil
}

func (this *Gofig) Unmarshall(dest any) error {
	root := reflect.ValueOf(dest)

	return this.applyRecords(root.Elem(), "")
}
