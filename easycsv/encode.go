package easycsv

import (
	"fmt"
	"reflect"
	"strconv"
)

func createDefaultConverterFromType(t reflect.Type) reflect.Value {
	switch t.Kind() {
	case reflect.Int:
		return reflect.ValueOf(strconv.Atoi)
	case reflect.Int8:
		return reflect.ValueOf(func(s string) (int8, error) {
			i, err := strconv.ParseInt(s, 0, 8)
			return int8(i), err
		})
	case reflect.Int16:
		return reflect.ValueOf(func(s string) (int16, error) {
			i, err := strconv.ParseInt(s, 0, 16)
			return int16(i), err
		})
	case reflect.Int32:
		return reflect.ValueOf(func(s string) (int32, error) {
			i, err := strconv.ParseInt(s, 0, 32)
			return int32(i), err
		})
	case reflect.Int64:
		return reflect.ValueOf(func(s string) (int64, error) {
			i, err := strconv.ParseInt(s, 0, 64)
			return int64(i), err
		})
	case reflect.Uint8:
		return reflect.ValueOf(func(s string) (uint8, error) {
			i, err := strconv.ParseUint(s, 0, 8)
			return uint8(i), err
		})
	case reflect.Uint16:
		return reflect.ValueOf(func(s string) (uint16, error) {
			i, err := strconv.ParseUint(s, 0, 16)
			return uint16(i), err
		})
	case reflect.Uint32:
		return reflect.ValueOf(func(s string) (uint32, error) {
			i, err := strconv.ParseUint(s, 0, 32)
			return uint32(i), err
		})
	case reflect.Uint64:
		return reflect.ValueOf(func(s string) (uint64, error) {
			i, err := strconv.ParseUint(s, 0, 32)
			return uint64(i), err
		})
	case reflect.Float32:
		return reflect.ValueOf(func(s string) (float32, error) {
			f, err := strconv.ParseFloat(s, 32)
			return float32(f), err
		})
	case reflect.Float64:
		return reflect.ValueOf(func(s string) (float64, error) {
			f, err := strconv.ParseFloat(s, 64)
			return float64(f), err
		})
	case reflect.Bool:
		return reflect.ValueOf(strconv.ParseBool)
	case reflect.String:
		return reflect.ValueOf(func(s string) (string, error) {
			return s, nil
		})
	default:
		return reflect.ValueOf(nil)
	}
}

func createDefaultConverter(field reflect.StructField, enc string) (reflect.Value, error) {
	c := createDefaultConverterFromType(field.Type)
	var err error
	if !c.IsValid() {
		err = fmt.Errorf("Unexpected field type for %s: %s", field.Name, field.Type)
	}
	return c, err
}
