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

	if !c.Next() {
		// вначале всегда идёт название плагина. Парсеру
		return fmt.Errorf("got no config data for plugin")
	}
	pluginName := c.Val()

	stream := newStream(c)

	err := unmarshal(stream, destValue.Elem())
	if err != nil {
		return err
	}

	if stream.Next() {
		return locErrf(stream.Token(), "got unexpected data for plugin '%s'", pluginName)
	}

	return nil
}

func unmarshal(s Stream, v reflect.Value) error {
	referenceType, isJSONUnmarshaler := refType(v.Type())
	if isJSONUnmarshaler {
		return processJSONUnmarshaler(s, v)
	}

	switch referenceType.Kind() {
	case reflect.Bool:
		return processBoolean(s, v)
	case reflect.Int8:
		return processInt8(s, v)
	case reflect.Int16:
		return processInt16(s, v)
	case reflect.Int32:
		return processInt32(s, v)
	case reflect.Int64:
		return processInt64(s, v)
	case reflect.Uint8:
		return processUint8(s, v)
	case reflect.Uint16:
		return processUint16(s, v)
	case reflect.Uint32:
		return processUint32(s, v)
	case reflect.Uint64:
		return processUint64(s, v)
	case reflect.Float32:
		return processFloat32(s, v)
	case reflect.Float64:
		return processFloat64(s, v)
	case reflect.String:
		return processString(s, v)
	case reflect.Slice:
	case reflect.Map:
	case reflect.Struct:
	}
	return nil
}

func processString(s Stream, v reflect.Value) error {
	if err := needArgValue(s, v); err != nil {
		return err
	}

	t := s.Token()
	r := ref(v)
	r.Set(reflect.ValueOf(t.Value))

	s.Confirm()
	return nil
}

func processBoolean(s Stream, v reflect.Value) error {
	if err := needArgValue(s, v); err != nil {
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

func processJSONUnmarshaler(s Stream, v reflect.Value) error {
	// получаем токен, на JSONUnmarshaler-ы у нас сильное ограничение – они должны записываться в один токен
	if err := needArgValue(s, v); err != nil {
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

func needArgValue(s Stream, v reflect.Value) error {
	if !s.NextArg() {
		return fmt.Errorf("got no data for %s", v.Type())
	}
	return nil
}
