package gofig

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func (this *Gofig) applyValue(dest reflect.Value, value any) error {
	if value == nil {
		return nil
	}

	src := reflect.ValueOf(value)

	// If types match directly, just set
	if src.Type().AssignableTo(dest.Type()) {
		dest.Set(src)
		return nil
	}

	// If convertible (e.g. int32 -> int64), use reflect conversion
	if src.Type().ConvertibleTo(dest.Type()) {
		dest.Set(src.Convert(dest.Type()))
		return nil
	}

	// Otherwise, stringify src and parse into dest type
	str, err := this.toString(src)
	if err != nil {
		return fmt.Errorf("cannot convert %T to string: %w", value, err)
	}

	return this.parseString(dest, str)
}

func (this *Gofig) toString(v reflect.Value) (string, error) {
	// time.Time
	if t, ok := v.Interface().(time.Time); ok {
		return t.Format(this.TimeFormats[0]), nil
	}

	switch v.Kind() {
	case reflect.String:
		return v.String(), nil
	case reflect.Bool:
		return strconv.FormatBool(v.Bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// Handle time.Duration specially
		if v.Type() == reflect.TypeOf(time.Duration(0)) {
			return time.Duration(v.Int()).String(), nil
		}
		return strconv.FormatInt(v.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64), nil
	}

	return "", fmt.Errorf("unsupported source kind: %s", v.Kind())
}

func (this *Gofig) parseString(dest reflect.Value, s string) error {
	// time.Duration
	if dest.Type() == reflect.TypeFor[time.Duration]() {
		d, err := time.ParseDuration(s)
		if err != nil {
			return err
		}
		dest.SetInt(int64(d))
		return nil
	}

	// time.Time
	if dest.Type() == reflect.TypeFor[time.Time]() {

		// try parse each
		for _, format := range this.TimeFormats {
			t, err := time.Parse(format, s)

			if err != nil {
				continue
			}

			dest.Set(reflect.ValueOf(t))
			return nil
		}

		return fmt.Errorf("Unable to parse date, valid formats: %v", this.TimeFormats)
	}

	switch dest.Kind() {
	case reflect.String:
		dest.SetString(s)

	case reflect.Bool:
		v, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		dest.SetBool(v)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		dest.SetInt(v)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		dest.SetUint(v)

	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		dest.SetFloat(v)

	default:
		return fmt.Errorf("unsupported dest kind: %s", dest.Kind())
	}

	return nil
}
