package caddycfg

import (
	"encoding/json"
	"fmt"
	"github.com/mholt/caddy"
	"reflect"
)

// Unmarshal unmarshaller into dest, which must not be channel
func Unmarshal(c *caddy.Controller, dest interface{}) error {
	destValue := reflect.ValueOf(dest)

	if destValue.Type().Kind() != reflect.Ptr {
		return fmt.Errorf("unmarshal into non-pointer %T", dest)
	}

	stream := newStream(c)
	if !c.Next() {
		// вначале всегда идёт название плагина.
		return fmt.Errorf("got no config data for plugin")
	}
	unmarshaler := &caddyCfgUnmarshaler{
		headToken: stream.Token(),
	}
	stream.Confirm()

	err := unmarshaler.unmarshal(stream, destValue.Elem())
	if err != nil {
		return err
	}

	if stream.Next() {
		return locErrf(stream.Token(), "got unexpected data for plugin '%s'", unmarshaler.headToken)
	}

	return nil
}

type caddyCfgUnmarshaler struct {
	headToken Token
}

func (c *caddyCfgUnmarshaler) unmarshal(s Stream, v reflect.Value) error {
	referenceType, isJSONUnmarshaler := refType(v.Type())
	if isJSONUnmarshaler {
		return c.processJSONUnmarshaler(s, v)
	}

	switch referenceType.Kind() {
	case reflect.Bool:
		return c.processBoolean(s, v)
	case reflect.Int8:
		return c.processInt8(s, v)
	case reflect.Int16:
		return c.processInt16(s, v)
	case reflect.Int32:
		return c.processInt32(s, v)
	case reflect.Int64:
		return c.processInt64(s, v)
	case reflect.Int:
		return c.processInt(s, v)
	case reflect.Uint8:
		return c.processUint8(s, v)
	case reflect.Uint16:
		return c.processUint16(s, v)
	case reflect.Uint32:
		return c.processUint32(s, v)
	case reflect.Uint64:
		return c.processUint64(s, v)
	case reflect.Uint:
		return c.processUint(s, v)
	case reflect.Float32:
		return c.processFloat32(s, v)
	case reflect.Float64:
		return c.processFloat64(s, v)
	case reflect.String:
		return c.processString(s, v)
	case reflect.Slice:
		return c.processSlice(s, v)
	case reflect.Map:
		return c.processMap(s, v)
	case reflect.Struct:
	}
	return nil
}

func (c *caddyCfgUnmarshaler) processMap(s Stream, v reflect.Value) error {
	keyType := v.Type().Key()
	switch keyType.String() {
	case
		reflect.Bool.String(),
		reflect.Int8.String(),
		reflect.Int16.String(),
		reflect.Int32.String(),
		reflect.Int64.String(),
		reflect.Int.String(),
		reflect.Uint8.String(),
		reflect.Uint16.String(),
		reflect.Uint32.String(),
		reflect.Uint64.String(),
		reflect.Uint.String(),
		reflect.String.String():
	default:
		rt, _ := refType(v.Type())
		return fmt.Errorf("unmarshaling into a %s is not supported: key can only be one of integer number type, boolean and string", rt)
	}

	if !s.NextArg() {
		return locErrf(c.headToken, "{ expected")
	}
	if s.Token().Value != "{" {
		return locErrf(s.Token(), "{ was expected, got %s", s.Token())
	}
	prevToken := s.Token()
	s.Confirm()

	r := refValue(v)
	dest := reflect.Zero(r.Type())
	valueType := r.Type().Elem()
	var closed bool
	for s.Next() {
		t := s.Token()
		prevToken = t
		if t.Value == "}" {
			closed = true
			s.Confirm()
			break
		}

		key := reflect.New(keyType)
		if err := c.unmarshal(s, key.Elem()); err != nil {
			return err
		}
		value := reflect.New(valueType)
		if err := c.unmarshal(s, value.Elem()); err != nil {
			return err
		}
		if dest.IsNil() {
			dest = reflect.MakeMap(r.Type())
		}
		dest.SetMapIndex(key.Elem(), value.Elem())
	}
	r.Set(dest)
	if !closed {
		return locErrf(prevToken, "} expected")
	}
	return nil
}

// There are two choices:
// 1. Slice of primitive types (JSONUnmarshaler, boolean, numeric types, string) can be represented in two ways
//     a) root a₁ a₂ … aₙ
//     b) root {
//           a₁
//           a₂
//           …
//           aₙ
//        }
// 2. Slice of complex types can only be represented as b variant
func (c *caddyCfgUnmarshaler) processSlice(s Stream, v reflect.Value) error {
	s.NextArg()
	if s.Token().Value == "{" {
		return c.processBlockedSlice(s, v)
	}

	r := ref(v)
	l := reflect.Zero(r.Type())
	for s.NextArg() {
		if s.Token().Value == "{" {
			rt, _ := refType(v.Type())
			return locErrf(s.Token(), "unmarshal block with arguments into %s", rt)
		}
		sliceElementType := l.Type().Elem()
		sliceItem := reflect.New(sliceElementType)
		rr := sliceItem.Elem()
		if err := c.unmarshal(s, rr); err != nil {
			return err
		}
		l = reflect.Append(l, rr)
	}
	r.Set(l)

	return nil
}

func (c *caddyCfgUnmarshaler) processBlockedSlice(s Stream, v reflect.Value) error {
	prevToken := s.Token()
	s.Confirm() // we reached { to be in here, so passing it

	r := ref(v)
	l := reflect.Zero(r.Type())

	// read until closing }
	var closed bool
	for s.Next() {
		t := s.Token()
		prevToken = t
		if t.Value == "}" {
			closed = true
			s.Confirm()
			break
		}
		sliceElementType := l.Type().Elem()
		sliceItem := reflect.New(sliceElementType)
		rr := sliceItem.Elem()
		if err := c.unmarshal(s, rr); err != nil {
			return err
		}
		l = reflect.Append(l, rr)
	}
	if !closed {
		return locErrf(prevToken, "} expected")
	}
	r.Set(l)

	return nil
}

func (c *caddyCfgUnmarshaler) processString(s Stream, v reflect.Value) error {
	if err := c.needArgValue(s, v); err != nil {
		return err
	}

	t := s.Token()
	r := ref(v)
	r.Set(reflect.ValueOf(t.Value))

	s.Confirm()
	return nil
}

func (c *caddyCfgUnmarshaler) processBoolean(s Stream, v reflect.Value) error {
	if err := c.needArgValue(s, v); err != nil {
		return err
	}

	t := s.Token()
	r := ref(v)
	switch t.Value {
	case "true":
		r.Set(reflect.ValueOf(true))
	case "false":
		r.Set(reflect.ValueOf(false))
	default:
		return locErrf(t, "true or false expected, got %s", t)
	}

	s.Confirm()

	return nil
}

func (c *caddyCfgUnmarshaler) processJSONUnmarshaler(s Stream, v reflect.Value) error {
	// получаем токен, на JSONUnmarshaler-ы у нас сильное ограничение – они должны записываться в один токен
	if err := c.needArgValue(s, v); err != nil {
		return err
	}

	t := s.Token()
	r := ref(v).Interface().(json.Unmarshaler)
	if err := r.UnmarshalJSON([]byte(t.Value)); err != nil {
		if err = r.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, t.Value))); err != nil {
			return locErrf(t, "cannot unmarshal: %s", err)
		}
	}

	s.Confirm()

	return nil
}

func (c *caddyCfgUnmarshaler) needArgValue(s Stream, v reflect.Value) error {
	if !s.NextArg() {
		return locErrf(c.headToken, "got no data for %s", v.Type())
	}
	return nil
}
