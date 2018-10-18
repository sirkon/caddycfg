package caddycfg

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type tokenError struct {
	Token
	err error
}

// Error ...
func (te tokenError) Error() string {
	return fmt.Sprintf("%s:%d: %s", te.File, te.Lin, te.err)
}

// TokenError diagnose a error caused with a given token
func TokenError(t Token, err error) error {
	return tokenError{
		Token: t,
		err:   err,
	}
}

// TokenErrorf diagnose a error caused with a given token with custom error message
func TokenErrorf(t Token, format string, a ...interface{}) error {
	err := fmt.Errorf(format, a...)
	return TokenError(t, err)
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

func createStructIndex(index map[string][]int, v reflect.Value, prefix []int) error {
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)

		if field.Type.Kind() == reflect.Struct && field.Anonymous {
			if err := createStructIndex(index, v.Field(i), append(prefix, i)); err != nil {
				return err
			}
			continue
		} else if len(field.Name) > 0 {
			letter := field.Name[:1]
			if letter == strings.ToLower(letter) {
				// avoid private fields
				continue
			}
		}

		name, ok := field.Tag.Lookup("json")
		if !ok {
			return fmt.Errorf("field '%s' from %s doesn't have 'json' tag", field.Name, v.Field(i).Type())
		}
		if _, ok := index[name]; ok {
			return fmt.Errorf("field '%s' from %s has duplicate json tag value '%s'", field.Name, v.Field(i).Type(), name)
		}
		ppp := make([]int, len(prefix), len(prefix)+1)
		copy(ppp, prefix)
		index[name] = append(ppp, i)
	}
	return nil
}

func orderFields(index map[string][]int) []string {
	names := make([]string, len(index))
	i := 0
	for name := range index {
		names[i] = name
		i++
	}
	sort.Slice(names, func(i, j int) bool {
		indexI := index[names[i]]
		indexJ := index[names[j]]
		for ii, valueI := range indexI {
			if valueI < indexJ[ii] {
				return true
			}
			if valueI > indexJ[ii] {
				return false
			}
		}
		return true
	})
	return names
}
