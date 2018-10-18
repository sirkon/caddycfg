package caddycfg

import (
	"reflect"
	"strconv"
)

func (c *caddyCfgUnmarshaler) processInt8(s Stream, v reflect.Value) error {
	if err := c.needArgValue(s, v); err != nil {
		return err
	}

	t := s.Token()
	r := ref(v)
	value, err := strconv.ParseInt(t.Value, 10, 8)
	if err != nil {
		return TokenError(t, err)
	}
	r.Set(reflect.ValueOf(int8(value)))

	s.Confirm()

	return nil
}

func (c *caddyCfgUnmarshaler) processInt16(s Stream, v reflect.Value) error {
	if err := c.needArgValue(s, v); err != nil {
		return err
	}

	t := s.Token()
	r := ref(v)
	value, err := strconv.ParseInt(t.Value, 10, 16)
	if err != nil {
		return TokenError(t, err)
	}
	r.Set(reflect.ValueOf(int16(value)))

	s.Confirm()

	return nil
}

func (c *caddyCfgUnmarshaler) processInt32(s Stream, v reflect.Value) error {
	if err := c.needArgValue(s, v); err != nil {
		return err
	}

	t := s.Token()
	r := ref(v)
	value, err := strconv.ParseInt(t.Value, 10, 32)
	if err != nil {
		return TokenError(t, err)
	}
	r.Set(reflect.ValueOf(int32(value)))

	s.Confirm()

	return nil
}

func (c *caddyCfgUnmarshaler) processInt64(s Stream, v reflect.Value) error {
	if err := c.needArgValue(s, v); err != nil {
		return err
	}

	t := s.Token()
	r := ref(v)
	value, err := strconv.ParseInt(t.Value, 10, 64)
	if err != nil {
		return TokenError(t, err)
	}
	r.Set(reflect.ValueOf(value))

	s.Confirm()

	return nil
}

func (c *caddyCfgUnmarshaler) processInt(s Stream, v reflect.Value) error {
	if err := c.needArgValue(s, v); err != nil {
		return err
	}

	t := s.Token()
	r := ref(v)
	value, err := strconv.Atoi(t.Value)
	if err != nil {
		return TokenError(t, err)
	}
	r.Set(reflect.ValueOf(value))

	s.Confirm()

	return nil
}

func (c *caddyCfgUnmarshaler) processUint8(s Stream, v reflect.Value) error {
	if err := c.needArgValue(s, v); err != nil {
		return err
	}

	t := s.Token()
	r := ref(v)
	value, err := strconv.ParseUint(t.Value, 10, 8)
	if err != nil {
		return TokenError(t, err)
	}
	r.Set(reflect.ValueOf(uint8(value)))

	s.Confirm()

	return nil
}

func (c *caddyCfgUnmarshaler) processUint16(s Stream, v reflect.Value) error {
	if err := c.needArgValue(s, v); err != nil {
		return err
	}

	t := s.Token()
	r := ref(v)
	value, err := strconv.ParseUint(t.Value, 10, 16)
	if err != nil {
		return TokenError(t, err)
	}
	r.Set(reflect.ValueOf(uint16(value)))

	s.Confirm()

	return nil
}

func (c *caddyCfgUnmarshaler) processUint32(s Stream, v reflect.Value) error {
	if err := c.needArgValue(s, v); err != nil {
		return err
	}

	t := s.Token()
	r := ref(v)
	value, err := strconv.ParseUint(t.Value, 10, 32)
	if err != nil {
		return TokenError(t, err)
	}
	r.Set(reflect.ValueOf(uint32(value)))

	s.Confirm()

	return nil
}

func (c *caddyCfgUnmarshaler) processUint64(s Stream, v reflect.Value) error {
	if err := c.needArgValue(s, v); err != nil {
		return err
	}

	t := s.Token()
	r := ref(v)
	value, err := strconv.ParseUint(t.Value, 10, 64)
	if err != nil {
		return TokenError(t, err)
	}
	r.Set(reflect.ValueOf(value))

	s.Confirm()

	return nil
}

func (c *caddyCfgUnmarshaler) processUint(s Stream, v reflect.Value) error {
	if err := c.needArgValue(s, v); err != nil {
		return err
	}

	t := s.Token()
	r := ref(v)
	value, err := strconv.ParseUint(t.Value, 10, 64)
	if err != nil {
		return TokenError(t, err)
	}
	r.Set(reflect.ValueOf(uint(value)))

	s.Confirm()

	return nil
}

func (c *caddyCfgUnmarshaler) processFloat32(s Stream, v reflect.Value) error {
	if err := c.needArgValue(s, v); err != nil {
		return err
	}

	t := s.Token()
	r := ref(v)
	value, err := strconv.ParseFloat(t.Value, 32)
	if err != nil {
		return TokenError(t, err)
	}
	r.Set(reflect.ValueOf(float32(value)))

	s.Confirm()

	return nil
}

func (c *caddyCfgUnmarshaler) processFloat64(s Stream, v reflect.Value) error {
	if err := c.needArgValue(s, v); err != nil {
		return err
	}

	t := s.Token()
	r := ref(v)
	value, err := strconv.ParseFloat(t.Value, 64)
	if err != nil {
		return TokenError(t, err)
	}
	r.Set(reflect.ValueOf(value))

	s.Confirm()

	return nil
}
