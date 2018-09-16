package caddycfg

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

type tmpType struct {
	a string
}

func (a *tmpType) UnmarshalJSON(data []byte) error {
	a.a = string(data)
	return nil
}

func TestReference(t *testing.T) {
	var dest ***int
	r := ref(reflect.ValueOf(&dest).Elem())
	r.Set(reflect.ValueOf(112))
	require.Equal(t, 112, ***dest)

	var input **tmpType
	r = ref(reflect.ValueOf(&input).Elem())
	r.Interface().(json.Unmarshaler).UnmarshalJSON([]byte("1234"))
	require.Equal(t, "1234", (*input).a)
}

func TestRefType(t *testing.T) {
	type ownTmpType struct {
		json.Unmarshaler
	}

	var dest ***ownTmpType
	whatType, isJSONUnmarshaler := refType(reflect.ValueOf(&dest).Elem().Type())
	require.True(t, isJSONUnmarshaler)
	require.IsType(t, &ownTmpType{}, reflect.Zero(whatType).Interface())
}

func TestLocerrf(t *testing.T) {
	tok := Token{
		File:  "caddyfile",
		Value: "abcdef",
		Lin:   1,
		Col:   1,
	}
	err := locErrf(tok, "error happened: %s", errors.New("error"))
	require.Equal(t, "caddyfile:1: error happened: error", err.Error())
}
