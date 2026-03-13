package gofig

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

const LABEL = "conf"

type Records = map[string]any

type Gofig struct {
	Fields        Fields
	Records       Records
	TimeFormats   []string
	MissingValues []string
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

func (this *Gofig) Read(provider Provider) error {

	src, err := provider.Source()

	if src == nil || err != nil {
		return err
	}

	for _, path := range this.Fields {
		value, err := src.Read(path)

		if err != nil {
			return err
		}

		if value != nil {
			this.Records[path] = value
		}
	}

	return nil
}

type NoValuesError struct {
	MissingValues []string
}

func (this *NoValuesError) Error() string {
	return fmt.Sprintf("Missing values: %v", this.MissingValues)
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
				this.MissingValues = append(this.MissingValues, currentPath)
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
	this.MissingValues = make([]string, 0)
	err := this.applyRecords(root.Elem(), "")

	if len(this.MissingValues) != 0 {
		return errors.Join(&NoValuesError{this.MissingValues}, err)
	}

	return err
}
