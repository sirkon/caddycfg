package caddycfg

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type tokenError struct {
	Token
	err error
}

// Error ...
func (te tokenError) Error() string {
	return fmt.Sprintf("%s:%d: %s", te.File, te.Lin, te.err)
}

func locErr(t Token, err error) error {
	return tokenError{
		Token: t,
		err:   err,
	}
}

func locErrf(t Token, format string, a ...interface{}) error {
	err := fmt.Errorf(format, a...)
	return locErr(t, err)
}

func ref(v reflect.Value) reflect.Value {
	return refValue(v)
}

func refValue(v reflect.Value) reflect.Value {
	if v.Type().Kind() == reflect.Ptr {
		nv := reflect.New(v.Type().Elem())
		v.Set(nv)
		if _, isJSONUnmarshaler := v.Interface().(json.Unmarshaler); isJSONUnmarshaler {
			return v
		}
		return refValue(nv.Elem())
	}
	return v
}

func refType(t reflect.Type) (reflect.Type, bool) {
	if _, isJSONUnmarshaler := reflect.Zero(t).Interface().(json.Unmarshaler); isJSONUnmarshaler {
		return t, true
	}
	if t.Kind() == reflect.Ptr {
		return refType(t.Elem())
	}
	return t, false
}
